package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// These type definitions make it possible to
// end every type with "boi"
type IntyBoi int
type FloatyBoi float64
type StringyBoi string

func boiError(boiInputs ...interface{}) {
	fmt.Print("\033[31;1mBoi! ")
	fmt.Print(boiInputs...)
	fmt.Println(", boi\033[0m")
}

func main() {
	boiArgs := os.Args[1:] // boi
	if len(boiArgs) < 1 {
		boiError("Usage: boi script.boi\n")
		os.Exit(1)
	}
	err := boiBoi(boiArgs[0]) // boi
	if err != nil {
		boiError(err)
	}
}

func boiBoi(boiFilename string) error {
	if boiFilename[len(boiFilename)-3:] != "boi" {
		return fmt.Errorf(
			"boi %s: MUST end with 'boi'", boiFilename,
		)
	}
	boiFile, err := os.Open(boiFilename)
	if err != nil {
		return err
	}

	code, err := ioutil.ReadAll(boiFile)
	if err != nil {
		return err
	}

	lex := NewBoiInterpreter(code)
	if err := lex.Run(); err != nil {
		return err
	}

	return nil
}

type BoiVar struct {
	data []byte
}

const (
	BoiTokenValue = 1 // A string
	BoiTokenVar   = 2
	BoiTokenBoi   = 3 // End of statement
)

const (
	// BoiStateStatement means we're expecting a statement
	BoiStateStatement IntyBoi = 0 // boi
)

type Token struct {
	BoiType  IntyBoi
	BoiValue []byte
}

type BoiFunc interface {
	Do(args []BoiVar) error
}

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

type BoiContext struct {
	functions map[string]BoiFunc
	variables map[string]BoiVar
	parentCtx *BoiContext
}

func (ctx *BoiContext) Call(fname string, args []BoiVar) error {
	f, exists := ctx.functions[fname]
	if !exists {
		if ctx.parentCtx == nil {
			return fmt.Errorf("call to undefined function %s", fname)
		} else {
			return ctx.parentCtx.Call(fname, args)
		}
	}
	return f.Do(args)
}

type BoiInterpreter struct {
	input []byte
	pos   IntyBoi
	state IntyBoi

	rIsBoiVar *regexp.Regexp
	rIsBoi    *regexp.Regexp

	context *BoiContext
}

func NewBoiInterpreter(input []byte) *BoiInterpreter {
	rootContext := &BoiContext{
		map[string]BoiFunc{},
		map[string]BoiVar{},
		nil,
	}

	boi := &BoiInterpreter{
		input, 0, BoiStateStatement,
		nil, nil,
		rootContext,
	}
	boi.rIsBoiVar = regexp.MustCompile("^boi:[A-z][A-z0-9]*")
	boi.rIsBoi = regexp.MustCompile("^boi[\\s\\n]")

	// Add internal functions
	boi.context.functions["say"] = BoiFuncSay{}
	boi.context.functions["set"] = BoiFuncSet{boi}

	return boi
}

func (boi *BoiInterpreter) Run() error {
	for {
		if boi.whitespace() {
			return nil
		}
		if err := boi.doStatement(); err != nil {
			return err
		}
	}
	return nil
}

func (boi *BoiInterpreter) whitespace() bool {
	if !(boi.pos < IntyBoi(len(boi.input)-1)) {
		return true // reached EOF
	}
	for boi.input[boi.pos] == ' ' ||
		boi.input[boi.pos] == '\n' ||
		boi.input[boi.pos] == '\t' {
		boi.pos++
	}
	return false
}

func (boi *BoiInterpreter) noeof(hasEof bool) error {
	if hasEof {
		return errors.New("unexpected EOF")
	}
	return nil
}

func (boi *BoiInterpreter) doStatement() error {
	stmt, err := boi.getStatement()
	if err != nil {
		return err
	}
	return boi.ExecStmt(stmt)
}

func (boi *BoiInterpreter) getStatement() (*BoiStatement, error) {
	op := string(boi.input[boi.pos : boi.pos+4])
	switch op {
	case "boi!":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens := []Token{}
		tokBytes := []BoiVar{}
		for {
			boi.noeof(boi.whitespace())
			if token, err := boi.eatToken(); err == nil {
				if token.BoiType != BoiTokenBoi {
					tokens = append(tokens, token)
					tokBytes = append(tokBytes, BoiVar{token.BoiValue})
				} else {
					break
				}
			} else {
				return nil, err
			}
		}
		return &BoiStatement{
			BoiOpCall, tokens,
		}, nil
	default:
		return nil, errors.New("unexpected")
	}
	return nil, errors.New("unexpected")
}

func (boi *BoiInterpreter) Call(identifier string, args []BoiVar) error {
	if f, exists := boi.context.functions[identifier]; exists {
		err := f.Do(args)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("function %s: not found", identifier)
	}
	return nil
}

func (boi *BoiInterpreter) getValueOf(tok Token) (BoiVar, bool) {
	switch tok.BoiType {
	case BoiTokenValue:
		return BoiVar{tok.BoiValue}, true
	case BoiTokenVar:
		identifier := string(tok.BoiValue)
		value, exists := boi.context.variables[identifier]
		if !exists {
			// TODO: Raise error if strictboi
		}
		return value, exists
	}
	return BoiVar{}, false
}

func (boi *BoiInterpreter) eatToken() (Token, error) {
	if !(boi.pos < IntyBoi(len(boi.input))) {
		return Token{}, errors.New("unexpected EOF")
	}

	isBoi := boi.rIsBoi.Match(boi.input[boi.pos:])
	if isBoi {
		boi.pos += 4
		t := Token{
			BoiType:  BoiTokenBoi,
			BoiValue: []byte{},
		}
		return t, nil
	}

	var tokType IntyBoi = BoiTokenValue

	isBoiVar := boi.rIsBoiVar.Match(boi.input[boi.pos:])
	if isBoiVar {
		boi.pos += 4
		tokType = BoiTokenVar
	}

	if boi.input[boi.pos] == '"' {
		boi.pos++ // otherwise we'll stop at the first quote
		value := []byte{}
		literal := false
		for ; boi.pos < IntyBoi(len(boi.input)); boi.pos++ {
			c := boi.input[boi.pos]
			if literal {
				value = append(value, c)
			} else {
				if c == '\\' {
					literal = true
				} else if c == '"' {
					boi.pos++ // don't forget to go past this quote
					break
				} else {
					value = append(value, c)
				}
			}
		}
		t := Token{
			BoiType:  tokType,
			BoiValue: value,
		}
		return t, nil
	}
	if true {
		value := []byte{}
		literal := false
		for ; boi.pos < IntyBoi(len(boi.input)); boi.pos++ {
			c := boi.input[boi.pos]
			if literal {
				value = append(value, c)
			} else {
				if c == '\\' {
					literal = true
				} else if c == ' ' {
					boi.pos++ // don't forget to go past this space
					break
				} else {
					value = append(value, c)
				}
			}
		}
		t := Token{
			BoiType:  tokType,
			BoiValue: value,
		}
		return t, nil
	}
	return Token{}, errors.New("unexpected token")
}
