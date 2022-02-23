package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ecnepsnai/discord"
	"github.com/kqzz/mcgo"
)

var errAccIgnored error = errors.New("account was ignored, either commented or otherwise")

// from gosnipe lmao
func pingMojang() (float64, error) {
	var sumNanos int64
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		return 0, err
	}

	defer conn.Close()
	for i := 0; i < 3; i++ {
		recv := make([]byte, 4096)
		time1 := time.Now()
		conn.Write([]byte("PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(recv)
		sumNanos += time.Since(time1).Nanoseconds()
	}

	sumNanos /= 3
	return float64(sumNanos / 1000000), nil
}

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

	if strSliceContainsMultiOption(strSplitLower, []string{"mj"}) {
		AccType = mcgo.Mj
	} else if strSliceContainsMultiOption(strSplitLower, []string{"prename", "msprename", "msaprename", "pr"}) {
		AccType = mcgo.MsPr
	} else {
		AccType = mcgo.Ms
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
			case 3:
				{
					account = mcgo.MCaccount{
						Email:    strSplit[0],
						Password: strSplit[1],
						Type:     AccType,
					}
				}
			case 6:
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
			if len(strSplit) == 3 {
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
	return strings.ReplaceAll(fmt.Sprintf("%-9s", timestamp.Format("05.999999")), " ", "0")
}


func formatAccount(account *mcgo.MCaccount) string {
	return fmt.Sprintf("%v:%v | Type: %v", account.Email, account.Password, prettyAccType(account.Type))
}

func announceSnipe(username, auth string, account *mcgo.MCaccount) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.mcsniperpy.com/announce?username=%v&prename=%v", username, account.Type == mcgo.MsPr), nil)

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
	default:
		return errors.New(fmt.Sprintf("Got: %v Couldnt announce your snipe, please contact staff to manually post it.", res.StatusCode))
	}
}

func customServerAnnounce(name string) error {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	discord.WebhookURL = config.Announce.WebhookURL

	discord.Post(discord.PostOptions{
		Embeds: []discord.Embed{
			{
				Footer: &discord.Footer{
					Text:    "MCsniperGO",
					IconURL: "https://cdn.discordapp.com/icons/734794891258757160/a_011d19e6e17a5eb46d108fd45b28dc9d.webp?size=96",
				},
				Title:       "Successful Snipe",
				URL:         "https://github.com/Kqzz/MCsniperGO",
				Color:       3118847,
				Description: fmt.Sprintf("Name: [`%v`](https://namemc.com/search?q=%v)", name, name),
			},
		},
	})
	return nil
}

func accID(acc *mcgo.MCaccount) string {
	if acc.Email != "" {
		return acc.Email
	} else if acc.Bearer != "" {
		return acc.Bearer[len(acc.Bearer)-10:]
	} else {
		return "<unknown account>"
	}
}

func authAccount(acc *mcgo.MCaccount, droptime time.Time) error {
	// authenticating if bearer isn't loaded
	if acc.Bearer == "" {
		switch acc.Type {
		case mcgo.MsPr, mcgo.Ms:
			{
				var err error
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

		expAt, err := BearerExpiresAt(acc.Bearer)	

		if err != nil {
			log("warn", "failed to validate bearer: %v", err)
		}

		if droptime.After(expAt) {
			expires := "expires"
			if time.Now().After(expAt) {
				expires = "expired"
			}
			return fmt.Errorf("bearer %v at %v, before droptime", expires, expAt)
		} else {
			log("info", "bearer expires in: %v", time.Until(expAt))
		}
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

func prettyStatus(status int) string {
	color := "red"
	if status < 300 {
		color = "green"
	}
	return fmt.Sprintf("<fg=%v;op=underscore>%v</>", color, status)
}

func BearerExpiresAt(bearer string) (time.Time, error) {

	s := strings.Split(bearer, ".")
	if len(s) < 2 {
		return time.Time{}, errors.New("bearer not formatted properly")
	}

	decoded, err := base64.RawStdEncoding.DecodeString(s[1])
	if err != nil {
		return time.Time{}, err
	}

	type Token struct {
		Exp int64 `json:"exp"`
	}

	var token Token
	err = json.Unmarshal(decoded, &token)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(token.Exp, 0), nil
}
