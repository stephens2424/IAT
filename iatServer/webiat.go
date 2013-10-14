package main

import (
  "math"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"stephensearles.com/iat"
	"strconv"
)

var currentSubject = 0
var currentSubjects = make(map[int]*iat.ReadyExperiment)

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
		readyExperiment := experiment.MakeFrames()
		renderedFrames := make([]RenderedFrame, len(readyExperiment.Frames))
		frameBuffer := bytes.NewBuffer(make([]byte, 0, 200))
		for i, frame := range readyExperiment.Frames {
			err := frame.RenderFrame(frameBuffer, tmpl)
			if err != nil {
				continue
			}
			renderedFrames[i] = RenderedFrame{frameBuffer.String(), frame.Correct()}
			frameBuffer.Reset()
		}
    currentSubject += 1
    readyExperiment.Subject = currentSubject
		frameJson, err := json.Marshal(renderedExperiment{
			renderedFrames,
      currentSubject,
		})
		if err == nil {
			framesBuffer <- frameJson
      currentSubjects[currentSubject] = &readyExperiment
		}
	}
}

func receiveData(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subjectIDString := req.FormValue("subjectID")
  subjectID, err := strconv.Atoi(subjectIDString)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  currentExperiment, ok := currentSubjects[subjectID]
  if !ok {
    http.Error(w, "Invalid subjectID", http.StatusBadRequest)
    return
  }

	times := req.Form["times[]"]
  rawResponses := req.Form["resp[]"]
  responses := make([]iat.Response, 0, len(times))

  prev := 0.0
	for i, time := range times {
		f, err := strconv.ParseFloat(time, 64)
		if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
    var d iat.Direction
    if rawResponses[i] == "l" {
      d = iat.Left
    } else {
      d = iat.Right
    }
    responses = append(responses, iat.Response{Time:f - prev, Dir: d})
    prev = f
	}
  score := currentExperiment.CalculateScore(responses)
  fmt.Fprintf(w, "Score: %v\n",score)
  absScore := math.Abs(score)
  if absScore < .15 {
    fmt.Fprintf(w, "You show no significant bias")
    return
  } else if absScore < .35 {
    fmt.Fprintf(w, "You show mild bias")
  } else if absScore < .65 {
    fmt.Fprintf(w, "You show some bias")
  } else {
    fmt.Fprintf(w, "You show significant bias")
  }
  fmt.Fprintf(w, " associating ")
  if score > 0 {
    fmt.Fprintf(w, "%s with %s and %s with %s",
      currentExperiment.DichotomyA.ListA.Title,
      currentExperiment.DichotomyB.ListA.Title,
      currentExperiment.DichotomyA.ListB.Title,
      currentExperiment.DichotomyB.ListB.Title,
    )
  } else {
    fmt.Fprintf(w, "%s with %s and %s with %s",
      currentExperiment.DichotomyA.ListA.Title,
      currentExperiment.DichotomyB.ListB.Title,
      currentExperiment.DichotomyA.ListB.Title,
      currentExperiment.DichotomyB.ListA.Title,
    )
  }
}

type renderedExperiment struct {
	Frames []RenderedFrame
  SubjectID int
}

type RenderedFrame struct {
	HTML    string
	Correct iat.Direction
}
