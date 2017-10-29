package main

import (
	"errors"
	"fmt"
)

// BoiGoFunc defines a function that can be invoked
// by the Boi-lang script
type BoiGoFunc func(*BoiContext, []BoiVar) (BoiVar, error)

// BoiGoFuncScruct defines an interface that can be
// invoked as a function by the Boi-lang script
type BoiGoFuncStruct interface {
	Run(*BoiContext, []BoiVar) (BoiVar, error)
}

// BoiGoFunctionAdapter implements the function interface
// that the Boi-lang interpreter requires in order to make
// it possible to invoke an implementor of the BoiGoFunc
// interface
type BoiGoFunctionAdapter struct {
	function    BoiGoFuncStruct
	interpreter *BoiInterpreter
}

func (adapter BoiGoFunctionAdapter) Do(args []BoiVar) error {
	context := adapter.interpreter.subContext()
	defer adapter.interpreter.returnContext()
	returnValue, err := adapter.function.Run(context, args)
	if err != nil {
		return err
	}
	adapter.interpreter.context.variables["exit"] = returnValue
	return nil
}

// BoiGoFuncAsFuncStruct makes it possible to pass a function
// matching the BoiGoFunc type wherever one might pass an
// implementor of BoiGoFuncStruct (since this is what
// BoiGoFunctionAdapter aggregates)
type BoiGoFuncAsFuncStruct struct {
	function BoiGoFunc
}

func (structure BoiGoFuncAsFuncStruct) Run(
	ctx *BoiContext, args []BoiVar,
) (
	BoiVar, error,
) {
	return structure.function(ctx, args)
}

/*
func BoiFuncName(context *BoiContext, args []BoiVar) ([]byte, error) {
	return nil, nil
}
*/

func BoiFuncSay(context *BoiContext, args []BoiVar) (BoiVar, error) {
	for _, bvar := range args {
		fmt.Print(string(bvar.data))
	}
	fmt.Println()
	return BoiVar{}, nil
}

func BoiFuncSet(context *BoiContext, args []BoiVar) (BoiVar, error) {
	if len(args) < 2 {
		return BoiVar{}, errors.New("set requires 2 parameters")
	}
	key := string(args[0].data)
	//context.parentCtx.variables[key] = args[1]
	context.parentCtx.Set(key, args[1])
	return args[1], nil
}

func BoiFuncCat(context *BoiContext, args []BoiVar) (BoiVar, error) {
	output := []byte{}
	for _, arg := range args {
		output = append(output, arg.data...)
	}
	return BoiVar{output}, nil
}
