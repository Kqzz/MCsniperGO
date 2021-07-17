package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kqzz/mcgo"
)

var commentedError error = errors.New("acc is commented out")

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func loadAccStr(accStr string) (mcgo.MCaccount, error) {
	if strings.HasPrefix(accStr, "#") {
		return mcgo.MCaccount{}, commentedError
	}
	var account mcgo.MCaccount
	strSplit := strings.Split(accStr, ":")
	switch len(strSplit) {
	case 2:
		{
			account = mcgo.MCaccount{
				Email:             strSplit[0],
				Password:          strSplit[1],
				SecurityQuestions: []mcgo.SqAnswer{},
				SecurityAnswers:   []string{},
				Bearer:            "",
				UUID:              "",
				Username:          "",
			}
		}
	case 5:
		{
			account = mcgo.MCaccount{
				Email:             strSplit[0],
				Password:          strSplit[1],
				SecurityQuestions: []mcgo.SqAnswer{},
				SecurityAnswers:   strSplit[2:5],
				Bearer:            "",
				UUID:              "",
				Username:          "",
			}
		}
	default:
		{
			return account, errors.New("wrong number of values, needs to be formatted email:password or email:password:answer:answer:answer")
		}
	}
	return account, nil
}

func loadAccSlice(accSlice []string) []mcgo.MCaccount {
	var accounts []mcgo.MCaccount
	for i, accStr := range accSlice {
		acc, err := loadAccStr(accStr)
		if err != nil {
			if !errors.Is(err, commentedError) {
				logErr(fmt.Sprintf(`got error "%v" while loading acc on line %v\n`, err, i+1))
			}
			continue
		}
		accounts = append(accounts, acc)
	}
	return accounts
}

func genHeader() string {
	header := `
███╗   ███╗ ██████╗███████╗███╗   ██╗██╗██████╗ ███████╗██████╗  ██████╗  ██████╗ 
████╗ ████║██╔════╝██╔════╝████╗  ██║██║██╔══██╗██╔════╝██╔══██╗██╔════╝ ██╔═══██╗
██╔████╔██║██║     ███████╗██╔██╗ ██║██║██████╔╝█████╗  ██████╔╝██║  ███╗██║   ██║
██║╚██╔╝██║██║     ╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  ██╔══██╗██║   ██║██║   ██║
██║ ╚═╝ ██║╚██████╗███████║██║ ╚████║██║██║     ███████╗██║  ██║╚██████╔╝╚██████╔╝
╚═╝     ╚═╝ ╚═════╝╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ 
<fg=white>Created by Kqzz#0001</>
`

	for _, char := range []string{"╗", "║", "╝", "╔", "═"} {
		header = strings.ReplaceAll(header, char, fmt.Sprintf("<fg=white>%v</>", char))
	}

	header = strings.ReplaceAll(header, "█", "<fg=cyan>█</>")

	return header
}
