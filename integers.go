package main

import (
	"math/big"
	"strconv"
)

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

type BoiFuncSub struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncSub) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	sum := new(big.Int)
	sum = sum.SetUint64(0)
	for _, arg := range args {
		tmp := new(big.Int)
		tmp = tmp.SetBytes(arg.data)
		sum = sum.Sub(sum, tmp)
	}

	context.variables["exit"] = BoiVar{sum.Bytes()}
	return nil
}

type BoiFuncDiv struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncDiv) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	sum := new(big.Int)
	sum = sum.SetUint64(0)
	for _, arg := range args {
		tmp := new(big.Int)
		tmp = tmp.SetBytes(arg.data)
		sum = sum.Div(sum, tmp)
	}

	context.variables["exit"] = BoiVar{sum.Bytes()}
	return nil
}

type BoiFuncMul struct {
	interpreter *BoiInterpreter
}

func (f BoiFuncMul) Do(args []BoiVar) error {
	context := f.interpreter.subContext()
	defer f.interpreter.returnContext()

	sum := new(big.Int)
	sum = sum.SetUint64(0)
	for _, arg := range args {
		tmp := new(big.Int)
		tmp = tmp.SetBytes(arg.data)
		sum = sum.Mul(sum, tmp)
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
