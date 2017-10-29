package main

import (
	"fmt"
)

const (
	BoiOpCall = 1
	BoiOpIf   = 2
	BoiOpLoop = 3
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
		identifier := string(stmt.Tokens[0].BoiValue)

		args := []BoiVar{}
		for _, tok := range stmt.Tokens[1:] {
			value, _ := boi.getValueOf(tok)
			args = append(args, value)
		}

		return boi.Call(identifier, args)
	case BoiOpIf:
		if len(stmt.Tokens) < 1 {
			return fmt.Errorf("boi? must have at least one token")
		}
		identifier := string(stmt.Tokens[0].BoiValue)

		args := []BoiVar{}
		for _, tok := range stmt.Tokens[1:] {
			value, _ := boi.getValueOf(tok)
			args = append(args, value)
		}

		// Call statement
		err := boi.Call(identifier, args)
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
		identifier := string(stmt.Tokens[0].BoiValue)
		continueLoop := true

		for continueLoop {

			// Recalculate arguments
			args := []BoiVar{}
			for _, tok := range stmt.Tokens[1:] {
				value, _ := boi.getValueOf(tok)
				args = append(args, value)
			}

			// Call statement
			err := boi.Call(identifier, args)
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
	}
	return fmt.Errorf(
		"internal error (aka boi is broken): invalid op code: %d",
		stmt.Operation,
	)
}
