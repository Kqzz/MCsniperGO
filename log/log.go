package log

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/Kqzz/MCsniperGO/pkg/mc"
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

func Log(level, message string, params ...interface{}) {
	format, e := formats[level]
	if !e {
		format = "%s"
	}

	color.Printf(format, fmt.Sprintf(message, params...))
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
    // Replace Windows line endings with Unix line endings
    unixLineEndings := strings.ReplaceAll(headerTxt, "\r\n", "\n")
    
    // Split the string, handling both single and double line breaks
    headers := strings.Split(unixLineEndings, "\n\n")
    
    // Remove any empty entries
    var nonEmptyHeaders []string
    for _, header := range headers {
        if strings.TrimSpace(header) != "" {
            nonEmptyHeaders = append(nonEmptyHeaders, header)
        }
    }
    
    src := rand.New(rand.NewSource(time.Now().UnixNano()))
    i := src.Intn(len(nonEmptyHeaders))
	i := src.Intn(len(headers) - 1)
	return fmt.Sprintf("\033[38;5;8m%v\033[0m\n\n<fg=blue;op=bold>MCsniperGO</> - made by kqzz (kqzz.me)\n\n", headers[i])
}

func FmtTimestamp(timestamp time.Time) string {
	return strings.ReplaceAll(fmt.Sprintf("%-9s", timestamp.Format("05.999999")), " ", "0")
}

func PrettyStatus(status int) string {
	color := "red"
	if status < 300 && status > 199 {
		color = "green"
	}
	return fmt.Sprintf("<fg=%v;op=underscore>%v</>", color, status)
}

func GetDropRange() mc.DropRange {
	for {
		rawDroptimes := Input("droptime range (start-end/infinite)")

		if rawDroptimes == "inf" || rawDroptimes == "infinite" {
			return mc.DropRange{Start: time.Time{}, End: time.Time{}}
		}

		rawDroptimesSplit := strings.Split(rawDroptimes, "-")

		if len(rawDroptimesSplit) != 2 {
			Log("err", "invalid droptime range")
			continue
		}

		startDroptimeNum, err := strconv.Atoi(rawDroptimesSplit[0])
		if err != nil {
			Log("err", "invalid droptime start")
			continue
		}
		endDroptimeNum, err := strconv.Atoi(rawDroptimesSplit[1])
		if err != nil {
			Log("err", "invalid droptime end")
			continue
		}
		startDroptime := time.Unix(int64(startDroptimeNum), 0)
		endDroptime := time.Unix(int64(endDroptimeNum), 0)

		return mc.DropRange{Start: startDroptime, End: endDroptime}
	}

}

func LastQuarter(s string) string {
	length := len(s)
	quarter := length / 4
	return s[length-quarter:]
}
