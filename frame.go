package iat

import (
	"html/template"
	"io"
)

type frame interface {
	RenderFrame(wr io.Writer, tmpl *template.Template) error
	Correct() Direction
  Block() int
}

type iatFrame struct {
	Center      string
	correct     Direction
  block int
	FrameTitles frameTitles
}

func (f iatFrame) RenderFrame(wr io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(wr, f)
}

func (f iatFrame) Correct() Direction {
	return f.correct
}

func (f iatFrame) Block() int {
  return f.block
}

type frameTitles struct {
	UpperLeft, LowerLeft, UpperRight, LowerRight string
}

func singleDichotomyFrames(a Dichotomy, trials int, leftIsA bool, block int) []frame {
	frames := make([]frame, trials)
	combinedItems := NewRandomLeftRightList([][]string{a.ListA.Items}, [][]string{a.ListB.Items})

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
		item, dir := combinedItems.Get()
		frames[i] = iatFrame{item, dir, block, titles}
	}
	return frames
}

func doubleDichotomyFrames(a, b Dichotomy, trials int, leftIsAUpper bool, block int) []frame {
	frames := make([]frame, trials)
	var leftLists, rightLists [][]string

	var titles frameTitles
	if leftIsAUpper {
		titles = frameTitles{
			UpperLeft:  a.ListA.Title,
			UpperRight: a.ListB.Title,
		}
		leftLists = [][]string{a.ListA.Items, b.ListA.Items}
		rightLists = [][]string{a.ListB.Items, b.ListB.Items}
	} else {
		titles = frameTitles{
			UpperLeft:  a.ListB.Title,
			UpperRight: a.ListA.Title,
		}
		leftLists = [][]string{a.ListB.Items, b.ListA.Items}
		rightLists = [][]string{a.ListA.Items, b.ListB.Items}
	}
	combinedItems := NewRandomLeftRightList(leftLists, rightLists)
	titles.LowerLeft = b.ListA.Title
	titles.LowerRight = b.ListB.Title
	for i := 0; i < trials; i++ {
		item, dir := combinedItems.Get()
		frames[i] = iatFrame{item, dir, block, titles}
	}
	return frames
}
