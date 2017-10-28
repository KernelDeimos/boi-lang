package main

import (
	"errors"
	"fmt"
)

type BoiFuncSay struct{}

func (f BoiFuncSay) Do(args []BoiVar) error {
	for _, bvar := range args {
		fmt.Print(string(bvar.data))
	}
	fmt.Println()
	return nil
}

type BoiFuncSet struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncSet) Do(args []BoiVar) error {
	if len(args) < 2 {
		return errors.New("set requires 2 parameters")
	}
	key := string(args[0].data)
	f.interpreter.context.variables[key] = args[1]
	return nil
}

type BoiFuncCat struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncCat) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	if len(args) < 2 {
		return errors.New("cat requires 2 parameters")
	}
	output := []byte{}
	for _, arg := range args {
		output = append(output, arg.data...)
	}
	context.variables["exit"] = BoiVar{output}
	return nil
}
