package main

import (
	"errors"
	"math/big"
	"math/rand"
	"strconv"
	"time"
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

func BoiFuncLess(context *BoiContext, args []BoiVar) (BoiVar, error) {
	if len(args) < 2 {
		return BoiVar{}, errors.New("< requires at least 2 parameters")
	}
	for i := 0; i < len(args)-1; i++ {
		a, b := new(big.Int), new(big.Int)
		a = a.SetBytes(args[i].data)
		b = b.SetBytes(args[i+1].data)
		if a.Cmp(b) >= 0 {
			return BoiVar{[]byte("false")}, nil
		}
	}
	return BoiVar{[]byte("true")}, nil
}

func BoiFuncIsEven(context *BoiContext, args []BoiVar) (BoiVar, error) {
	if len(args) != 1 {
		return BoiVar{}, errors.New("IsEven can only take one value (for now)")
	}
	// always seed random
	rand.Seed(time.Now().UTC().UnixNano())
	probabilityOfWrongAnswer := rand.Intn(100)

	value := new(big.Int)
	value = value.SetBytes(args[0].data)

	valueInt := value.Uint64()

	even := valueInt%2 == 0

	if rand.Intn(100) < probabilityOfWrongAnswer {
		even = !even
	}

	probabilityOfEven := probabilityOfWrongAnswer
	if even {
		probabilityOfEven = 100 - probabilityOfWrongAnswer
	}

	return BoiVar{[]byte{byte(probabilityOfEven)}}, nil

}
