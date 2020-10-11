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
}

func (s Step) String() string {
	stringString := ""
	stringString += fmt.Sprintf("%s: Needs %d\nCook at %s\n", s.StepType, s.TimeNeeded, s.Temperature.String())
	stringString += s.Instructions + "\n"

	return stringString
}

func (st StepType) String() string {
	//this will panic if you try to pass in something not
	//in the constant array
	return stepTypeNames[st]

}
