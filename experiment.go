// Package IAT provides a framework to run an Implicit Association Test, described
// by Greenwald, Nosek, Banaji in the Journal of Personality and Social Psychology,
// 2003. At the time of this writing, the paper is available at:
//
// http://faculty.washington.edu/agg/pdf/GB&N.JPSP.2003.pdf
package iat

import (
	"crypto/rand"
	"math/big"
)

type Experiment struct {
	DichotomyA, DichotomyB Dichotomy
}

// Dichotomy represents a set of opposed categories for the experiment.
type Dichotomy struct {
	ListA, ListB CategoryList
}

// CategoryList represents the actual category underneath the dichotomy.
// This contains the title and words that will be displayed on the screen.
type CategoryList struct {
	Title string
	Items []string
}

// Function MakeFrames builds the frames for the experiment. See Table 1 of Greenwald 2003.
func (e *Experiment) MakeFrames() []frame {
	frames := make([]frame, 0, 180)

	flip := randBool()

	frames = append(frames, singleDichotomyFrames(e.DichotomyA, 2, flip)...)
	frames = append(frames, singleDichotomyFrames(e.DichotomyB, 2, true)...) // Greenwald block 2 does not flip
	frames = append(frames, doubleDichotomyFrames(e.DichotomyA, e.DichotomyB, 2, flip)...)
	frames = append(frames, doubleDichotomyFrames(e.DichotomyA, e.DichotomyB, 4, flip)...)
	frames = append(frames, singleDichotomyFrames(e.DichotomyA, 2, !flip)...)
	frames = append(frames, doubleDichotomyFrames(e.DichotomyA, e.DichotomyB, 2, !flip)...)
	frames = append(frames, doubleDichotomyFrames(e.DichotomyA, e.DichotomyB, 4, !flip)...)

	return frames
}

type direction int

const (
	Left  direction = iota
	Right           = iota
)

type RandomLeftRightList struct {
	leftList, rightList []string
}

func NewRandomList(leftLists, rightLists []string) RandomList {
	this := RandomLeftRightList{}
	this.leftList = make([]string, 0, 10)
	this.rightList = make([]string, 0, 10)
	for _, l := range leftLists {
		this.leftList = append(this.leftLists, l...)
	}
	for _, l := range rightLists {
		this.rightList = append(this.rightLists, l...)
	}
	return this
}

func (r RandomList) Get() string {
	i, _ := randInt(uint64(len(r)))
	return r[i]
}

func randInt(max uint64) (r uint64, err error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return i.Uint64(), err
}

func randBool() bool {
	i, _ := randInt(1)
	if i == 0 {
		return false
	}
	return true
}
