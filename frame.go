package iat

import (
	"html/template"
	"io"
)

type frame interface {
	RenderFrame(wr io.Writer, tmpl *template.Template) error
}

type iatFrame struct {
	Center      string
	FrameTitles frameTitles
	Correct     string
}

func (f iatFrame) RenderFrame(wr io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(wr, f)
}

type frameTitles struct {
	UpperLeft, LowerLeft, UpperRight, LowerRight string
}

func singleDichotomyFrames(a Dichotomy, trials int, leftIsA bool) []frame {
	frames := make([]frame, trials)
	combinedItems := NewRandomList(a.ListA.Items, a.ListB.Items)

	var titles frameTitles
	if leftIsA {
		titles = frameTitles{
			UpperLeft:  a.ListA.Title,
			UpperRight: a.ListB.Title,
		}
	} else {
		titles = frameTitles{
			UpperLeft:  a.ListB.Title,
			UpperRight: a.ListA.Title,
		}
	}
	for i := 0; i < trials; i++ {
		frames[i] = iatFrame{combinedItems.Get(), titles}
	}
	return frames
}

func doubleDichotomyFrames(a, b Dichotomy, trials int, leftIsAUpper bool) []frame {
	frames := make([]frame, trials)
	combinedItems := NewRandomList(a.ListA.Items, a.ListB.Items, b.ListA.Items, b.ListB.Items)

	var titles frameTitles
	if leftIsAUpper {
		titles = frameTitles{
			UpperLeft:  a.ListA.Title,
			UpperRight: a.ListB.Title,
		}
	} else {
		titles = frameTitles{
			UpperLeft:  a.ListB.Title,
			UpperRight: a.ListA.Title,
		}
	}
	titles.LowerLeft = b.ListA.Title
	titles.LowerRight = b.ListB.Title
	for i := 0; i < trials; i++ {
		frames[i] = iatFrame{combinedItems.Get(), titles}
	}
	return frames
}
