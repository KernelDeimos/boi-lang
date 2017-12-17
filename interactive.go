package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func boiSlackServer(hostname string) {

	var lastContext *BoiContext = nil

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		text := c.PostForm("text")
		lex := NewBoiInterpreter([]byte(text))
		if lastContext != nil {
			lex.context = lastContext
		}

		// === + Terrible Hack Start ===
		runBoi := func() (string, error) {

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			outc := make(chan string)
			go func() {
				var buff bytes.Buffer
				io.Copy(&buff, r)
				outc <- buff.String()
			}()

			err := lex.Run()

			w.Close()
			os.Stdout = old

			return <-outc, err
		}

		output, err := runBoi()

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, struct {
				ResponseType string `json:"response_type"`
				Text         string `json:"text"`
			}{
				"in_channel", output,
			})
		}
		lastContext = lex.context
	})
	r.Run(hostname)
}
