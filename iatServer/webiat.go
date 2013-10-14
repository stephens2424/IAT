package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"stephensearles.com/iat"
	"strconv"
)

func main() {
	var experiment iat.Experiment
	f, err := os.Open("testExperiment.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	d := json.NewDecoder(f)
	err = d.Decode(&experiment)
	if err != nil {
		panic(err)
	}

	// Buffering 5 clients/users/experiments to handle a fairly low server load
	framesBuffer := make(chan []byte, 5)
	go bufferExperiments(experiment, framesBuffer)
	http.HandleFunc("/expdata/sample", func(w http.ResponseWriter, req *http.Request) {
		w.Write(<-framesBuffer)
	})

	http.HandleFunc("/postResults", http.HandlerFunc(receiveData))
	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.ListenAndServe(":36300", nil)
}

// bufferExperiments precalculates (randomizes) experiments since that component
// can be a tad slower than I'd like right now.
func bufferExperiments(experiment iat.Experiment, framesBuffer chan []byte) {
	tmpl := template.Must(template.ParseFiles("templates/frame.tmpl"))
	for {
		frames := experiment.MakeFrames()
		renderedFrames := make([]RenderedFrame, len(frames))
		frameBuffer := bytes.NewBuffer(make([]byte, 0, 200))
		for i, frame := range frames {
			err := frame.RenderFrame(frameBuffer, tmpl)
			if err != nil {
				continue
			}
			renderedFrames[i] = RenderedFrame{frameBuffer.String(), frame.Correct()}
			frameBuffer.Reset()
		}
		frameJson, err := json.Marshal(renderedExperiment{
			renderedFrames,
		})
		if err == nil {
			framesBuffer <- frameJson
		}
	}
}

func receiveData(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sum := 0.0
	prev := 0.0
	times := req.Form["times[]"]
	for _, time := range times {
		f, err := strconv.ParseFloat(time, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sum += f - prev
		prev = f
	}
	fmt.Fprintf(w, "times: %v\n", times)
	fmt.Fprintf(w, "Average response time %v", sum/float64(len(times)))
}

type renderedExperiment struct {
	Frames []RenderedFrame
}

type RenderedFrame struct {
	HTML    string
	Correct iat.Direction
}
