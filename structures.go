package main

import "fmt"

const (
	BoiOpCall = 1
)

type BoiStatement struct {
	Operation int
	Tokens    []Token
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
	}
	return fmt.Errorf(
		"internal error (aka boi is broken): invalid op code: %d",
		stmt.Operation,
	)
}
