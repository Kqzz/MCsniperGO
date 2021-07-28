package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kqzz/mcgo"
)

var errCommented error = errors.New("acc is commented out")

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
		return mcgo.MCaccount{}, errCommented
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
					}
				}
			case 5:
				{
					account = mcgo.MCaccount{
						Email:           strSplit[0],
						Password:        strSplit[1],
						SecurityAnswers: strSplit[2:4],
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
			}
		}

	}

	// switch len(strSplit) {
	// case 2:
	// 	{
	// 		if strSplit[1] == "bearer" || strSplit[1] == "br" {
	// 			var accType mcgo.AccType = mcgo.Mj
	// 			for _, v := range strSplit {
	// 				if strings.ToLower(v) == "ms" {
	// 					accType = mcgo.Ms
	// 				}
	// 			}
	// 			account = mcgo.MCaccount{
	// 				Bearer: strSplit[0],
	// 				Type:   accType,
	// 			}
	// 		} else {
	// 			account = mcgo.MCaccount{
	// 				Email:    strSplit[0],
	// 				Password: strSplit[1],
	// 				Type:     mcgo.Mj,
	// 			}
	// 		}
	// 	}
	// case 5:
	// 	{
	// 		account = mcgo.MCaccount{
	// 			Email:           strSplit[0],
	// 			Password:        strSplit[1],
	// 			SecurityAnswers: strSplit[2:5],
	// 			Type:            mcgo.Mj,
	// 		}
	// 	}
	// case 3, 4:
	// 	{
	// 		if strings.ToLower(strSplit[2]) == "ms" {
	// 			var prename bool = false
	// 			for _, v := range strSplit {
	// 				v = strings.ToLower(v)
	// 				if v == "prename" || v == "pr" {
	// 					prename = true
	// 				}
	// 			}
	// 			var accType mcgo.AccType
	// 			if prename {
	// 				accType = mcgo.MsPr
	// 			} else {
	// 				accType = mcgo.Ms
	// 			}
	// 			account = mcgo.MCaccount{
	// 				Email:    strSplit[0],
	// 				Password: strSplit[1],
	// 				Type:     accType,
	// 			}
	// 		} else {
	// 			return account, errors.New("wrong number of values, needs to be formatted email:password or email:password:answer:answer:answer or, for ms acc, email:password:ms (dont replace ms with anything)")
	// 		}
	// 	}
	// default:
	// 	{
	// 		return account, errors.New("wrong number of values, needs to be formatted email:password or email:password:answer:answer:answer")
	// 	}
	// }
	return account, nil
}

func loadAccSlice(accSlice []string) []*mcgo.MCaccount {
	var accounts []*mcgo.MCaccount
	for i, accStr := range accSlice {
		acc, err := loadAccStr(accStr)
		if err != nil {
			if !errors.Is(err, errCommented) {
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
