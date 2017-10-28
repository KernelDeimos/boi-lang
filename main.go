package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/davecgh/go-spew/spew"
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

const (
	BoiTokenValue = 1 // A string
	BoiTokenBoi   = 2 // End of statement
)

const (
	// BoiStateStatement means we're expecting a statement
	BoiStateStatement IntyBoi = 0 // boi
)

type Token struct {
	BoiType  IntyBoi
	BoiValue StringyBoi
}

type BoiFunc interface {
	Do(args []string)
}

type BoiInterpreter struct {
	input []byte
	pos   IntyBoi
	state IntyBoi

	rIsBoiVar *regexp.Regexp
	rIsBoi    *regexp.Regexp

	functions map[string]BoiFunc
}

func NewBoiInterpreter(input []byte) *BoiInterpreter {
	boi := &BoiInterpreter{
		input, 0, BoiStateStatement,
		nil, nil,
		map[string]BoiFunc{},
	}
	boi.rIsBoiVar = regexp.MustCompile("^boi:[A-z][A-z0-9]*")
	boi.rIsBoi = regexp.MustCompile("^boi[\\s\\n]")
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
	op := string(boi.input[boi.pos : boi.pos+4])
	switch op {
	case "boi!":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		identifier, err := boi.eatIdentifier()
		tokens := []Token{}
		for {
			boi.noeof(boi.whitespace())
			if token, err := boi.eatToken(); err == nil {
				spew.Dump(token)
				if token.BoiType == BoiTokenValue {
					tokens = append(tokens, token)
				} else {
					break
				}
			} else {
				return err
			}
		}
		if err != nil {
			return err
		}
		fmt.Printf("Want to call function %s with arguments:\n", identifier)
		for _, tok := range tokens {
			fmt.Printf("\t%d: %s\n", tok.BoiType, tok.BoiValue)
		}
		fmt.Printf("-----\n")
		return nil
	default:
		return errors.New("unexpected")
	}
	return errors.New("unexpected")
}

func (boi *BoiInterpreter) eatIdentifier() (string, error) {
	start := boi.pos
	end := boi.pos

	c := boi.input[boi.pos]
	if (c > 0x40 && c <= 0x5A) ||
		(c > 0x60 && c < 0x7A) {
		end++
	} else {
		return "", fmt.Errorf("char %d: invalid identifier", boi.pos)
	}

	for ; ; end++ {
		c := boi.input[end]
		if (c > 0x40 && c <= 0x5A) ||
			(c > 0x60 && c < 0x7A) ||
			(c >= 0x30 && c < 0x3A) {
			end++
		} else {
			break
		}
	}

	boi.pos = end
	return string(boi.input[start:end]), nil
}

func (boi *BoiInterpreter) eatToken() (Token, error) {
	fmt.Println("|" + string(boi.input[boi.pos:len(boi.input)]) + "|")

	if !(boi.pos < IntyBoi(len(boi.input))) {
		return Token{}, errors.New("unexpected EOF")
	}

	isBoi := boi.rIsBoi.Match(boi.input[boi.pos:])
	if isBoi {
		boi.pos += 4
		t := Token{
			BoiType:  BoiTokenBoi,
			BoiValue: "",
		}
		return t, nil
	}

	isBoiVar := boi.rIsBoiVar.Match(boi.input[boi.pos:])
	if isBoiVar {
		boi.pos += 4
		identifier, err := boi.eatIdentifier()
		if err != nil {
			return Token{}, err
		}
		t := Token{
			BoiType: BoiTokenValue,
			BoiValue: StringyBoi(
				fmt.Sprintf("placeholder for '%s'", identifier),
			),
		}
		return t, nil
	}

	if boi.input[boi.pos] == '"' {
		value := ""
		literal := false
		for ; boi.pos < IntyBoi(len(boi.input)); boi.pos++ {
			c := boi.input[boi.pos]
			if literal {
				value = value + string([]byte{c})
			} else {
				if c == '\\' {
					literal = true
				} else if c == '"' {
					boi.pos++ // don't forget to go past this quote
					break
				} else {
					value = value + string([]byte{c})
				}
			}
		}
		t := Token{
			BoiType:  BoiTokenValue,
			BoiValue: StringyBoi(value),
		}
		return t, nil
	}
	if true {
		value := ""
		literal := false
		for ; boi.pos < IntyBoi(len(boi.input)); boi.pos++ {
			c := boi.input[boi.pos]
			if literal {
				value = value + string([]byte{c})
			} else {
				if c == '\\' {
					literal = true
				} else if c == ' ' {
					boi.pos++ // don't forget to go past this space
					break
				} else {
					value = value + string([]byte{c})
				}
			}
		}
		t := Token{
			BoiType:  BoiTokenValue,
			BoiValue: StringyBoi(value),
		}
		return t, nil
	}
	return Token{}, errors.New("unexpected token")
}
