package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Feel free to add your droptime API to this file, the more droptime APIS that are used, the more stable this sniper should be (theoretically)

var (
	allApis = []string{"api.coolkidmacho.com", "drops.peet.ws"}
)

type coolkidmachoRespStruct struct {
	Unix int64 `json:"UNIX"`
}

// grabs droptime from api.coolkidmacho.com
func coolkidmachoDroptime(username string) (time.Time, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.coolkidmacho.com/droptime/%v", username))

	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return time.Time{}, err
	}

	if resp.StatusCode < 300 {
		var coolkidmachoResponse coolkidmachoRespStruct
		err = json.Unmarshal(respBytes, &coolkidmachoResponse)
		if err != nil {
			return time.Time{}, err
		}

		return time.Unix(coolkidmachoResponse.Unix, 0), nil
	}

	return time.Time{}, fmt.Errorf("failed to grab droptime with status %v and body %v", resp.Status, string(respBytes))
}

// TODO: implement grabbing droptime from drops.peet.ws once that API is fixed

type peetResponse struct {
	Unix int64  `json:"UNIX,omitempty"`
	Err  string `json:"error,omitempty"`
}

func peetDroptime(username string) (time.Time, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://drops.peet.ws/droptime?name=%v", username), nil)
	if err != nil {
		return time.Time{}, err
	}

	req.Header.Set("User-Agent", "PiratSnipe")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return time.Time{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return time.Time{}, fmt.Errorf("got status %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return time.Time{}, err
	}

	var peetResp peetResponse

	err = json.Unmarshal(body, &peetResp)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(peetResp.Unix, 0), nil
}

func apiFromStr(name string) func(string) (time.Time, error) {
	var api func(string) (time.Time, error)
	switch strings.ToLower(name) {
	case "ckm", "coolkid", "coolkidmacho", "api.coolkidmacho.com":
		api = coolkidmachoDroptime
	case "peet", "peet.ws", "drops.peet.ws":
		api = peetDroptime
	default:
		api = coolkidmachoDroptime
	}
	return api
}

func getDroptime(username, preference string) (time.Time, error) {
	dropFuncs := []func(string) (time.Time, error){}
	var droptime time.Time
	var err error

	prefApi := apiFromStr(preference)

	dropFuncs = append(dropFuncs, prefApi)

	for _, api := range allApis {
		dropFuncs = append(dropFuncs, apiFromStr(api))
	}

	for _, grabDrop := range dropFuncs {
		droptime, err = grabDrop(username)
		if err != nil {
			logErr(fmt.Sprintf("got err \"%v\" while requesting droptime...", err))
			time.Sleep(time.Second * 1)
			logInfo("trying next droptime api...")
		} else {
			break
		}
	}

	if droptime.IsZero() {
		return time.Time{}, errors.New("all droptime APIs failed to grab droptime")
	}

	return droptime, nil
}
