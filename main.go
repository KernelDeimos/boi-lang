package main

import (
	"errors"
	"fmt"
	"io"
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

	var reader io.Reader

	if len(boiArgs) < 1 {
		boiInteractive()
	} else if boiArgs[0] == "/slack" {
		hostname := ""
		if len(boiArgs) > 1 {
			hostname = boiArgs[1]
		}
		boiSlackServer(hostname)
	} else if boiArgs[0] == "-" {
		reader = os.Stdout
	} else {
		//
		boiFilename := boiArgs[0]

		if boiFilename[len(boiFilename)-3:] != "boi" {
			boiError(fmt.Errorf(
				"boi %s: MUST end with 'boi'", boiFilename,
			))
		}
		var err error
		reader, err = os.Open(boiFilename)
		if err != nil {
			boiError(err)
		}

	}

	err := boiBoi(reader) // boi
	if err != nil {
		boiError(err)
	}
}

func boiBoi(reader io.Reader) error {
	code, err := ioutil.ReadAll(reader)
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
	BoiTokenCall  = 4
)

const (
	// BoiStateStatement means we're expecting a statement
	BoiStateStatement IntyBoi = 0 // boi
)

// Enumerated list of "source types"
const (
	BoiSourceLocal  = 1
	BoiSourceReturn = 2
)

type Token struct {
	BoiType IntyBoi

	// For strings or variable names
	BoiValue []byte

	// Source context (for variables)
	BoiSource int

	Children []Token
}

type BoiFunc interface {
	Do(args []BoiVar) error
}

type BoiContext struct {
	functions map[string]BoiFunc
	variables map[string]BoiVar
	parentCtx *BoiContext
	returnCtx *BoiContext
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

func (ctx *BoiContext) Set(vname string, value BoiVar) error {
	tryContext := ctx
	_, exists := tryContext.variables[vname]
	for !exists {
		if tryContext.parentCtx == nil {
			break
		} else {
			tryContext = tryContext.parentCtx
		}
		_, exists = tryContext.variables[vname]
	}
	if exists {
		tryContext.variables[vname] = value
	} else {
		ctx.variables[vname] = value
	}
	return nil
}

func (ctx *BoiContext) Get(vname string) (BoiVar, bool) {
	v, exists := ctx.variables[vname]
	if !exists {
		if ctx.parentCtx == nil {
			// TODO: Raise error if boi.context is strict context
			return BoiVar{}, false
		} else {
			return ctx.parentCtx.Get(vname)
		}
	}
	return v, true
}

type BoiInterpreter struct {
	input []byte
	pos   IntyBoi
	state IntyBoi

	rSyntaxToken *regexp.Regexp

	rIsBoiVar *regexp.Regexp
	rIsRetVar *regexp.Regexp
	rIsBoi    *regexp.Regexp

	context *BoiContext
}

func NewBoiInterpreter(input []byte) *BoiInterpreter {
	rootContext := &BoiContext{
		map[string]BoiFunc{},
		map[string]BoiVar{},
		nil, nil,
	}

	boi := &BoiInterpreter{
		input, 0, BoiStateStatement,
		nil, nil, nil, nil,
		rootContext,
	}
	boi.rIsBoiVar = regexp.MustCompile("^boi:[A-Za-z][A-Za-z0-9]*")
	boi.rIsRetVar = regexp.MustCompile("^ret:[A-Za-z][A-Za-z0-9]*")
	boi.rIsBoi = regexp.MustCompile("^boi[\\s\\n]")

	boi.rSyntaxToken = regexp.MustCompile(
		`^([A-Za-z]+[!,:\?]?|[\[\];]|--)`,
	)

	// Add internal functions
	boi.RegisterGoFunction("say", BoiFuncSay)
	boi.RegisterGoFunction("set", BoiFuncSet)
	boi.RegisterGoFunction("icanhas", BoiFuncGet)
	boi.RegisterGoFunction("nyan", BoiFuncCat)
	boi.context.functions["int"] = BoiFuncInt{boi}
	boi.context.functions["+"] = BoiFuncAdd{boi}
	boi.context.functions["-"] = BoiFuncSub{boi}
	boi.context.functions["/"] = BoiFuncDiv{boi}
	boi.context.functions["*"] = BoiFuncMul{boi}
	boi.context.functions["dec"] = BoiFuncDec{boi}
	boi.RegisterGoFunction("<", BoiFuncLess)

	// Grey area (memes, also practical)
	boi.RegisterGoFunction("declare", BoiFuncDeclare)

	// Memes
	boi.RegisterGoFunction("IsEven", BoiFuncIsEven)

	return boi
}

func (boi *BoiInterpreter) RegisterGoFunction(fname string, f BoiGoFunc) {
	adapter := BoiGoFunctionAdapter{
		BoiGoFuncAsFuncStruct{f}, boi,
	}
	boi.context.functions[fname] = adapter
}

func (boi *BoiInterpreter) RegisterGoFunctionStruct(
	fname string, f BoiGoFuncStruct,
) {
	adapter := BoiGoFunctionAdapter{
		f, boi,
	}
	boi.context.functions[fname] = adapter
}

func (boi *BoiInterpreter) subContext() *BoiContext {
	ctx := &BoiContext{
		map[string]BoiFunc{},
		map[string]BoiVar{},
		boi.context, nil,
	}
	boi.context = ctx
	return ctx
}

func (boi *BoiInterpreter) returnContext() error {
	returnCtx := boi.context
	boi.context = boi.context.parentCtx
	if boi.context == nil {
		return errors.New("returned to nil context")
	}
	boi.context.returnCtx = returnCtx
	return nil
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
	for ; boi.pos < IntyBoi(len(boi.input)); boi.pos++ {
		//
		if !(boi.input[boi.pos] == ' ' ||
			boi.input[boi.pos] == '\n' ||
			boi.input[boi.pos] == '\t') {
			return false
		}
	}
	return true
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
	if stmt == nil {
		return nil
	}
	return boi.ExecStmt(stmt)
}

func (boi *BoiInterpreter) getStatement() (*BoiStatement, error) {
	op := string(boi.rSyntaxToken.Find(boi.input[boi.pos:]))

	for op == "--" {
		boi.pos += 2
		for {
			boi.pos++
			if boi.pos >= IntyBoi(len(boi.input)) {
				break
			}
			if boi.input[boi.pos] == '\n' {
				break
			}
		}
		if boi.whitespace() {
			return nil, nil
		}
		op = string(boi.rSyntaxToken.Find(boi.input[boi.pos:]))
	}

	switch op {
	case "boi!":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		return &BoiStatement{
			BoiOpCall, tokens, nil,
		}, nil
	case "boi,":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		return NewCallStatement("say", tokens), nil

	case "boi:":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		return NewCallStatement("set", tokens), nil

	case "one":
		fallthrough
	case "ONE":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		return NewCallStatement("declare", tokens), nil

	case "boi?":
		boi.pos += 4
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		statements, err := boi.GetStatements()
		if err != nil {
			return nil, err
		}

		return &BoiStatement{
			BoiOpIf, tokens, statements,
		}, nil
	case "bloop":
		boi.pos += 5
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		statements, err := boi.GetStatements()
		if err != nil {
			return nil, err
		}

		return &BoiStatement{
			BoiOpLoop, tokens, statements,
		}, nil
	case "oh":
		fallthrough
	case "OH":
		boi.pos += 2
		boi.noeof(boi.whitespace())
		tokens, err := boi.GetTokens()
		if err != nil {
			return nil, err
		}

		statements, err := boi.GetStatements()
		if err != nil {
			return nil, err
		}

		return &BoiStatement{
			BoiOpFuncDef, tokens, statements,
		}, nil
	case "BOI":
		boi.pos += 3
		return nil, nil
	default:
		return nil, fmt.Errorf("unrecognized keyword '%s'", op)
	}
}

func (boi *BoiInterpreter) GetTokens() ([]Token, error) {
	tokens := []Token{}
	for {
		boi.noeof(boi.whitespace())
		if token, err := boi.eatToken(); err == nil {
			if token.BoiType != BoiTokenBoi {
				tokens = append(tokens, token)
			} else {
				break
			}
		} else {
			return tokens, err
		}
	}
	return tokens, nil
}

func (boi *BoiInterpreter) GetStatements() ([]*BoiStatement, error) {
	// Aggregate statements until we hit a nil statement ("BOI")
	statements := []*BoiStatement{}
	for {
		if boi.whitespace() {
			return nil, fmt.Errorf("end of file before BOI")
		}
		stmt, err := boi.getStatement()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			break
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (boi *BoiInterpreter) Call(identifier string, args []BoiVar) error {
	return boi.context.Call(identifier, args)
	/*
		if f, exists := boi.context.functions[identifier]; exists {
			err := f.Do(args)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("function %s: not found", identifier)
		}
		return nil
	*/
}

func (boi *BoiInterpreter) getValueOf(tok Token) (BoiVar, bool) {
	switch tok.BoiType {
	case BoiTokenValue:
		return BoiVar{tok.BoiValue}, true
	case BoiTokenVar:
		context := boi.context
		if tok.BoiSource == BoiSourceReturn {
			context = boi.context.returnCtx
		}

		if context == nil {
			// TODO: Raise error if boi.context is strict context
			return BoiVar{}, false
		}

		identifier := string(tok.BoiValue)
		value, exists := context.Get(identifier)
		if !exists {
			// TODO: Raise error if strict context
		}
		return value, exists
	case BoiTokenCall:
		identifier := string(tok.Children[0].BoiValue)

		args := []BoiVar{}
		for _, tok := range tok.Children[1:] {
			value, _ := boi.getValueOf(tok)
			args = append(args, value)
		}

		// Call statement
		err := boi.Call(identifier, args)
		if err != nil {
			return BoiVar{}, false // TODO: Raise error
		}

		output := boi.context.returnCtx.variables["exit"]
		return BoiVar(output), true
	}
	return BoiVar{}, false
}

func (boi *BoiInterpreter) eatToken() (Token, error) {
	if !(boi.pos < IntyBoi(len(boi.input))) {
		return Token{}, errors.New("unexpected EOF")
	}

	keyword := string(boi.rSyntaxToken.Find(boi.input[boi.pos:]))
	isBoi := keyword == "boi" || keyword == "]" || keyword == ";" ||
		keyword == "BOI"
	if isBoi {
		boi.pos += IntyBoi(len(keyword))
		t := Token{
			BoiType:  BoiTokenBoi,
			BoiValue: []byte{},
		}
		return t, nil
	}

	token := Token{
		BoiType:   BoiTokenValue,
		BoiValue:  []byte{},
		BoiSource: BoiSourceLocal,
	}

	isBoiVar := boi.rIsBoiVar.Match(boi.input[boi.pos:])
	isRetVar := boi.rIsRetVar.Match(boi.input[boi.pos:])
	if isBoiVar || isRetVar {
		boi.pos += 4
		token.BoiType = BoiTokenVar
	}
	if isRetVar {
		token.BoiSource = BoiSourceReturn
	}

	if boi.input[boi.pos] == '[' || boi.input[boi.pos] == '!' {
		boi.pos++
		if boi.whitespace() {
			return token, fmt.Errorf("end of file before BOI")
		}
		toks, err := boi.GetTokens()
		if err != nil {
			return token, err
		}
		token.BoiType = BoiTokenCall
		token.Children = toks

		return token, nil
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
		token.BoiValue = value
		return token, nil
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
				} else if c == ']' || c == ';' {
					break
				} else {
					value = append(value, c)
				}
			}
		}
		token.BoiValue = value
		return token, nil
	}
	return Token{}, errors.New("unexpected token")
}
