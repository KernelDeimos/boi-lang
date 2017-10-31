package main

import (
	"fmt"
)

const (
	BoiOpCall    = 1
	BoiOpIf      = 2
	BoiOpLoop    = 3
	BoiOpFuncDef = 4
)

type BoiStatement struct {
	Operation int
	Tokens    []Token
	Children  []*BoiStatement
}

func NewCallStatement(fname string, tokens []Token) *BoiStatement {
	functionToken := Token{
		BoiType:   BoiTokenValue,
		BoiValue:  []byte(fname),
		BoiSource: BoiSourceLocal,
	}

	tokens = append([]Token{functionToken}, tokens...)

	return &BoiStatement{
		BoiOpCall, tokens, nil,
	}
}

func (boi *BoiInterpreter) ExecStmt(stmt *BoiStatement) error {
	switch stmt.Operation {
	case BoiOpCall:
		if len(stmt.Tokens) < 1 {
			return fmt.Errorf("boi! must have at least one token")
		}

		args := []BoiVar{}
		for _, tok := range stmt.Tokens {
			value, _ := boi.getValueOf(tok)
			args = append(args, value)
		}

		identifier := string(args[0].data)
		return boi.Call(identifier, args[1:])
	case BoiOpIf:
		if len(stmt.Tokens) < 1 {
			return fmt.Errorf("boi? must have at least one token")
		}

		args := []BoiVar{}
		for _, tok := range stmt.Tokens {
			value, _ := boi.getValueOf(tok)
			args = append(args, value)
		}

		// Call statement
		identifier := string(args[0].data)
		err := boi.Call(identifier, args[1:])
		if err != nil {
			return err
		}

		// Execute subsequent statements if output is true
		exitVar, exists := boi.context.returnCtx.variables["exit"]

		// Check for falsy values
		switch true {
		case !exists:
			fallthrough
		case len(exitVar.data) == 0:
			fallthrough
		case string(exitVar.data) == "false":
			// Falsey value encountered - do not execute aggregate statements
			return nil
		}

		// Scope down for aggregate statements
		boi.subContext()
		defer boi.returnContext()

		// Execute aggregate statements
		for _, stmt := range stmt.Children {
			if stmt != nil {
				err := boi.ExecStmt(stmt)
				if err != nil {
					return err
				}
			}
		}
		return nil
	case BoiOpLoop:
		if len(stmt.Tokens) < 1 {
			return fmt.Errorf("bloop must have at least one token")
		}
		continueLoop := true

		for continueLoop {

			// Recalculate arguments
			args := []BoiVar{}
			for _, tok := range stmt.Tokens {
				value, _ := boi.getValueOf(tok)
				args = append(args, value)
			}

			// Call statement
			identifier := string(args[0].data)
			err := boi.Call(identifier, args[1:])
			if err != nil {
				return err
			}

			// Execute subsequent statements if output is true
			exitVar, exists := boi.context.returnCtx.variables["exit"]

			// Check for falsy values
			switch true {
			case !exists:
				fallthrough
			case len(exitVar.data) == 0:
				fallthrough
			case string(exitVar.data) == "false":
				// Falsey value encountered - do not execute aggregate statements
				continueLoop = false
			}

			if !continueLoop {
				break
			}

			// Scope down for aggregate statements
			err = func() error {
				boi.subContext()
				defer boi.returnContext()

				// Execute aggregate statements
				for _, stmt := range stmt.Children {
					if stmt != nil {
						err := boi.ExecStmt(stmt)
						if err != nil {
							return err
						}
					}
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
		return nil
	case BoiOpFuncDef:

		var identifier string
		if len(stmt.Tokens) < 1 {
			identifier = ""
		} else {
			identifierBoi, _ := boi.getValueOf(stmt.Tokens[0])
			identifier = string(identifierBoi.data)
		}

		boi.RegisterGoFunctionStruct(
			identifier,
			NewBoiStatementsFunction(stmt.Children, boi),
		)

		return nil
	}
	return fmt.Errorf(
		"internal error (aka boi is broken): invalid op code: %d",
		stmt.Operation,
	)
}
