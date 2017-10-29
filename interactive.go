package main

import (
	"bufio"
	"os"
)

func boiInteractive() {
	userin := bufio.NewReader(os.Stdin)

	var lastContext *BoiContext = nil

	for {
		text, err := userin.ReadBytes('\n')
		if err != nil {
			boiError(err)
			break
		}

		lex := NewBoiInterpreter(text)
		if lastContext != nil {
			lex.context = lastContext
		}
		if err := lex.Run(); err != nil {
			boiError(err)
			continue
		}
		lastContext = lex.context
	}
}
