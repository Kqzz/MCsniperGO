package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

func log(m string, t string, params ...interface{}) {
	color.Printf((logType[t] + " » " + "%s\n"), fmt.Sprintf(m, params...))
}

func userInput(m string, params ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)
	var out string
	color.Printf("<fg=white>[</><fg=cyan;op=bold>input</><fg=white>]</> %s » ", fmt.Sprintf(m, params...))
	out, _ = reader.ReadString('\n')
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	return out
}
