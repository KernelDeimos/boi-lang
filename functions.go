package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
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

// BoiStatementsFunction implements BoiGoFuncStruct, and runs aggregate
// BoiStatement objects in order.
type BoiStatementsFunction struct {
	statements  []*BoiStatement
	interpreter *BoiInterpreter
}

func NewBoiStatementsFunction(
	statements []*BoiStatement,
	interpreter *BoiInterpreter,
) *BoiStatementsFunction {
	return &BoiStatementsFunction{
		statements, interpreter,
	}
}

func (f *BoiStatementsFunction) Run(
	ctx *BoiContext, args []BoiVar,
) (BoiVar, error) {
	for i, val := range args {
		ctx.Set(string("arg."+strconv.Itoa(i)), val)
	}
	for _, stmt := range f.statements {
		err := f.interpreter.ExecStmt(stmt)
		if err != nil {
			return BoiVar{}, err
		}
	}
	exitValue, exists := ctx.variables["exit"]
	if !exists {
		exitValue = BoiVar{[]byte{}}
	}
	return exitValue, nil
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

// BoiFuncDeclare is similar to BoiFuncSet, but does not require a value
// parameter. It instead initializes the variable to a completely random
// value to **ensure** the application programmer can't make assumptions
// about the value. Adding a value parameter anyway is undefined behaviour.
func BoiFuncDeclare(context *BoiContext, args []BoiVar) (BoiVar, error) {
	if len(args) < 1 {
		return BoiVar{}, errors.New("one requires 1 parameters")
	}
	key := string(args[0].data)
	//context.parentCtx.variables[key] = args[1]

	value := make([]byte, 4)
	_, err := rand.Read(value)
	if err != nil {
		return BoiVar{}, err
	}

	// We can't use .Set() here, because that tries to find the variable
	// in parent scopes.
	context.parentCtx.variables[key] = BoiVar{value}
	return BoiVar{value}, nil
}

func BoiFuncCat(context *BoiContext, args []BoiVar) (BoiVar, error) {
	output := []byte{}
	for _, arg := range args {
		output = append(output, arg.data...)
	}
	return BoiVar{output}, nil
}
