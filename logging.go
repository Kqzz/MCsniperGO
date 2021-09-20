package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

func logInfo(m string, params ...interface{}) {
	color.Printf("<fg=white>[</><fg=cyan;op=bold>info</><fg=white>]</> » %s\n", fmt.Sprintf(m, params...))
}

func logSuccess(m string, params ...interface{}) {
	color.Printf("<fg=white>[</><fg=green;op=bold>success</><fg=white>]</> » %s\n", fmt.Sprintf(m, params...))
}

func logErr(m string, params ...interface{}) {
	color.Printf("<fg=white>[</><fg=red;op=bold>err</><fg=white>]</> » %s\n", fmt.Sprintf(m, params...))
}

func logWarn(m string, params ...interface{}) {
	color.Printf("<fg=white>[</><fg=yellow;op=bold>warn</><fg=white>]</> » %s\n", fmt.Sprintf(m, params...))
}

func logFatal(m string, params ...interface{}) {
	color.Printf("<fg=white>[</><fg=red;op=bold>fatal err</><fg=white>]</> » %s\n", fmt.Sprintf(m, params...))
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
