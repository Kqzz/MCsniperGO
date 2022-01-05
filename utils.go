package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
				return account, fmt.Errorf("invalid split count of %v on line: %v", len(strSplit), accStr)
			}
		}
	case mcgo.Ms, mcgo.MsPr:
		{
			if len(strSplit) == 2 {
				account = mcgo.MCaccount{
					Email:    strSplit[0],
					Password: strSplit[1],
					Type:     AccType,
				}
			} else {
				account = mcgo.MCaccount{
					Email:    strSplit[0],
					Password: "oauth2-external",
					Type:     AccType,
				}
			}
		}
	}

	return account, nil
}

func loadAccSlice(accSlice []string) []*mcgo.MCaccount {
	var accounts []*mcgo.MCaccount
	for i, accStr := range accSlice {
		if accStr == "" {
			continue
		}
		acc, err := loadAccStr(accStr)
		if err != nil {
			if !errors.Is(err, errAccIgnored) {
				log("error", `got error "%v" while loading acc on line %v`, err, i+1)
			}
			continue
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

func prettyAccType(acc mcgo.AccType) string {
	return map[mcgo.AccType]string{
		mcgo.Mj:   "mojang",
		mcgo.Ms:   "microsoft",
		mcgo.MsPr: "microsoft prename",
	}[acc] + " account"
}

func countAccounts(accounts []*mcgo.MCaccount) (int, int) {
	normCount := 0
	prenameCount := 0
	for _, acc := range accounts {
		switch acc.Type {
		case mcgo.Mj, mcgo.Ms:
			{
				normCount += 1
			}
		case mcgo.MsPr:
			{
				prenameCount += 1
			}
		}
	}
	return normCount, prenameCount
}

func fmtTimestamp(timestamp time.Time) string {
	return timestamp.Format("15:04:05.9999")
}

func formatAccount(account *mcgo.MCaccount) string {
	return fmt.Sprintf("%v:%v | Type: %v", account.Email, account.Password, prettyAccType(account.Type))
}

func announceSnipe(username, auth string, account *mcgo.MCaccount) error {
	prename := "false"
	if account.Type == mcgo.MsPr {
		prename = "true"
	}
	url := fmt.Sprintf("https://api.mcsniperpy.com/announce?username=%v&prename=%v", username, prename)

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", auth)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		{
			return errors.New("invalid auth code")
		}
	case 429:
		{
			return errors.New("too many requests (ask staff to manually announce)")
		}
	}

	if res.StatusCode != 204 {
		log("error", "got unknown status code while announcing snipe: %v", res.StatusCode)
	}

	return nil
}

func accID(acc *mcgo.MCaccount) string {
	if acc.Email != "" {
		return acc.Email
	}

	if acc.Bearer != "" {
		return acc.Bearer[len(acc.Bearer)-10:]
	}

	return "<unknown account>"
}

func authAccount(acc *mcgo.MCaccount) error {
	// authenticating if bearer isn't loaded
	if acc.Bearer == "" {
		switch acc.Type {
		case mcgo.MsPr, mcgo.Ms:
			{
				var err Error
				if acc.Password == "oauth2-external" {
					err = acc.InitAuthFlow()
				} else {
					err = acc.MicrosoftAuthenticate()
				}
				if err != nil {
					return err
				}
				log("info", "authenticating %s through ms auth", accID(acc))
			}
		case mcgo.Mj:
			{
				err := acc.MojangAuthenticate()
				if err != nil {
					return err
				}
				log("info", "authenticating %s through mojang auth", accID(acc))
			}
		}
	} else {
		log("info", "authing %s through manual bearer", accID(acc))
	}

	return nil
}

func accReadyToSnipe(acc *mcgo.MCaccount) (bool, error) {
	switch acc.Type {
	case mcgo.MsPr:
		{
			canCreateProfile, err := acc.HasGcApplied()

			if err != nil {
				return false, fmt.Errorf("prename acc: %v", err)
			}

			return canCreateProfile, nil
		}
	case mcgo.Mj, mcgo.Ms:
		{
			nameChangeInfo, err := acc.NameChangeInfo()

			if err != nil {
				return false, fmt.Errorf("mj / ms acc: %v", err)
			}

			return nameChangeInfo.Namechangeallowed, nil

		}
	}
	return false, nil
}

func estimatedProcess(send, recv time.Time) time.Time {
	return time.UnixMilli((send.UnixMilli() + recv.UnixMilli()) / 2)
}

func prettyStatus(status int) string {
	color := "red"
	if status < 300 {
		color = "green"
	}
	return fmt.Sprintf("<fg=%v;op=underscore>%v</>", color, status)
}
