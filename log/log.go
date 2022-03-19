package log

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/gookit/color"
)

type Logger struct {
	Color bool
}

var formats = map[string]string{
	/* "info":  "[<fg=blue>info</>] %s\n",
	"err":   "[<fg=red>erro</>] %s\n",
	"warn":  "[<fg=yellow>warn</>] %s\n",
	"input": "[<fg=blue>input</>] %s: ", */
	"info":    "[<fg=blue>*</>] %s\n",
	"err":     "[<fg=red>*</>] %s\n",
	"warn":    "[<fg=yellow>*</>] %s\n",
	"success": "[<fg=green>*</>] %s\n",
	"input":   "[<fg=blue>*</>] %s: ",
}

func Log(l, m string, params ...interface{}) {
	format, e := formats[l]
	if !e {
		format = "%s"
	}

	color.Printf(format, fmt.Sprintf(m, params...))
}

func Input(m string, params ...interface{}) string {
	scanner := bufio.NewScanner(os.Stdin)
	color.Printf(formats["input"], fmt.Sprintf(m, params...))

	scanner.Scan()

	t := scanner.Text()
	return t
}

//go:embed headers.txt
var headerTxt string

func GetHeader() string {

	headers := strings.Split(headerTxt, "\n\n")
	rand.Seed(time.Now().UnixNano())

	i := rand.Intn(len(headers) - 1)
	return fmt.Sprintf("\033[38;5;8m%v\033[0m\n\n<fg=blue>MCsniperGO</> - Made by kqzz (kqzz.me)\n\n", headers[i])
}

func FmtTimestamp(timestamp time.Time) string {
	return strings.ReplaceAll(fmt.Sprintf("%-9s", timestamp.Format("05.999999")), " ", "0")
}

func isBetween(n, x, y int) bool {
	return n >= x && n <= y
}

// show if late / early / close / success
func PrettyTimestampStatus(timestamp, droptime time.Time, status int) string {

	statusFormat := "<fg=white>[</><fg=%s>%s</><fg=white>]</>"
	if status < 300 && status > 199 {
		return fmt.Sprintf(statusFormat, "green", "success")
	}

	// handle really late / really early
	// tbh could be better, seeing as it doesn't handle cases where droptime isn't on the second (theoretically impossible)

	if timestamp.Before(droptime) {
		return fmt.Sprintf(statusFormat, "red", "early")
	}

	if timestamp.After(droptime.Add(time.Second * 1)) {
		return fmt.Sprintf(statusFormat, "red", "late")
	}

	ns := timestamp.Nanosecond()

	if isBetween(ns, .78e+7, 1.3e+7) {
		return fmt.Sprintf(statusFormat, "green", "close")
	}

	if ns < 1e+7 {
		return fmt.Sprintf(statusFormat, "red", "early")
	}

	return fmt.Sprintf(statusFormat, "red", "late")

}

func PrettyStatus(status int) string {
	color := "red"
	if status < 300 && status > 199 {
		color = "green"
	}
	return fmt.Sprintf("<fg=%v;op=underscore>%v</>", color, status)
}
