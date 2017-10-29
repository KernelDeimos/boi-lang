package main

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

	output := []byte{}
	for _, arg := range args {
		output = append(output, arg.data...)
	}
	context.variables["exit"] = BoiVar{output}
	return nil
}

type BoiFuncInt struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncInt) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	sum := new(big.Int)
	sum = sum.SetUint64(0)

	for _, arg := range args {
		data := arg.data
		value, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return err
		}
		tmp := new(big.Int)
		tmp.SetUint64(value)
		sum = sum.Add(sum, tmp)
	}

	context.variables["exit"] = BoiVar{sum.Bytes()}
	return nil
}

type BoiFuncAdd struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncAdd) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	sum := new(big.Int)
	sum = sum.SetUint64(0)
	for _, arg := range args {
		tmp := new(big.Int)
		tmp = tmp.SetBytes(arg.data)
		sum = sum.Add(sum, tmp)
	}

	context.variables["exit"] = BoiVar{sum.Bytes()}
	return nil
}

type BoiFuncDec struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncDec) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	value := new(big.Int)
	value = value.SetUint64(0)
	if len(args) > 0 {
		value = value.SetBytes(args[0].data)
	}

	output := []byte(value.String())

	context.variables["exit"] = BoiVar{output}
	return nil
}
