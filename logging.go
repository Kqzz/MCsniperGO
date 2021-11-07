package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

var logType = map[string]string{
	"success": "<fg=white>[</><fg=green;op=bold>success</><fg=white>]</>",
	"info":    "<fg=white>[</><fg=cyan;op=bold>info</><fg=white>]</>",
	"warn":    "<fg=white>[</><fg=yellow;op=bold>warn</><fg=white>]</>",
	"error":   "<fg=white>[</><fg=red;op=bold>err</><fg=white>]</>",
	"fatal":   "<fg=white>[</><fg=red;op=bold>fatal err</><fg=white>]</>",
}

func log(t, m string, params ...interface{}) {
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
