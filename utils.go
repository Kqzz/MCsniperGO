package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kqzz/mcgo"
)

var errAccIgnored error = errors.New("account was ignored, either commented or otherwise")

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

func strSliceContainsMultiOption(s []string, strs []string) bool {
	for _, str := range strs {
		for _, v := range s {
			if v == str {
				return true
			}
		}
	}

	return false
}

func loadAccStr(accStr string) (mcgo.MCaccount, error) {
	if strings.HasPrefix(accStr, "#") {
		return mcgo.MCaccount{}, errAccIgnored
	}
	var account mcgo.MCaccount
	strSplit := strings.Split(accStr, ":")
	strSplitLower := strings.Split(strings.ToLower(accStr), ":")

	var AccType mcgo.AccType

	if strSliceContainsMultiOption(strSplitLower, []string{"ms", "msa"}) {
		AccType = mcgo.Ms
	} else if strSliceContainsMultiOption(strSplitLower, []string{"prename", "msprename", "msaprename", "pr"}) {
		AccType = mcgo.MsPr
	} else {
		AccType = mcgo.Mj
	}

	if strSliceContainsMultiOption(strSplitLower, []string{"bearer"}) {
		account = mcgo.MCaccount{
			Bearer: strSplit[0],
			Type:   AccType,
		}
		return account, nil
	}

	switch AccType {
	case mcgo.Mj:
		{
			switch len(strSplit) {
			case 2:
				{
					account = mcgo.MCaccount{
						Email:    strSplit[0],
						Password: strSplit[1],
						Type:     AccType,
					}
				}
			case 5:
				{
					account = mcgo.MCaccount{
						Email:           strSplit[0],
						Password:        strSplit[1],
						SecurityAnswers: strSplit[2:5],
						Type:            AccType,
					}
				}
			default:
				return account, fmt.Errorf("invalid split count of %v", len(strSplit))
			}
		}
	case mcgo.Ms, mcgo.MsPr:
		{
			account = mcgo.MCaccount{
				Email:    strSplit[0],
				Password: strSplit[1],
				Type:     AccType,
			}
		}

	}

	return account, nil
}

func loadAccSlice(accSlice []string) []*mcgo.MCaccount {
	var accounts []*mcgo.MCaccount
	for i, accStr := range accSlice {
		acc, err := loadAccStr(accStr)
		if err != nil {
			if !errors.Is(err, errAccIgnored) {
				logErr(fmt.Sprintf(`got error "%v" while loading acc on line %v`, err, i+1))
				continue
			}
		}
		accounts = append(accounts, &acc)
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
<fg=cyan>https://mcsniperpy.com</>

`

	for _, char := range []string{"╗", "║", "╝", "╔", "═"} {
		header = strings.ReplaceAll(header, char, fmt.Sprintf("<fg=white>%v</>", char))
	}

	header = strings.ReplaceAll(header, "█", "<fg=blue>█</>")

	return header
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func censor(str string, amt int) string {
	out := []rune(str)
	for i := range out {
		if i >= amt && amt >= 0 {
			break
		}
		out[i] = '*'
	}
	return string(out)
}
