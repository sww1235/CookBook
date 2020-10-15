package main

import (
	"fmt"
	"time"
)

//StepType has 4 recognized values, prep, cook, wait
//and then anything else as other. This selects what duration in a Recipe
//TimeNeeded is added to.
type StepType int

//declare enum for Steptype
const (
	Prep StepType = iota // 0
	Cook
	Wait
	Other
)

var stepTypeNames = [...]string{
	"Prep",
	"Cook",
	"Wait",
	"Other",
}

type Step struct {
	ID           int // database id
	TimeNeeded   time.Duration
	StepType     StepType
	Temperature  Temperature
	Instructions string
	Order        int // order of step in recipe
}

func (s Step) String() string {
	stringString := ""
	stringString += fmt.Sprintf("%s Step %d: Needs %d\nCook at %s\n", s.StepType, s.Order, s.TimeNeeded, s.Temperature.String())
	stringString += s.Instructions + "\n"

	return stringString
}

func (st StepType) String() string {
	//this will panic if you try to pass in something not
	//in the constant array
	return stepTypeNames[st]

}

// ByName implements sort.Interface for []Step
// based on the Name field
// https://golang.org/pkg/sort/#pkg-overview
type ByOrderS []Step

func (s ByOrderS) Len() int           { return len(s) }
func (s ByOrderS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByOrderS) Less(i, j int) bool { return s[i].Order < s[j].Order }

// ByID implements sort.Interface for []Step
// based on the ID field
// https://golang.org/pkg/sort/#pkg-overview
type ByIDS []Step

func (s ByIDS) Len() int           { return len(s) }
func (s ByIDS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByIDS) Less(i, j int) bool { return s[i].ID < s[j].ID }
