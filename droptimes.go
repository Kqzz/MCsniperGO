package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Feel free to add your droptime API to this file, the more droptime APIS that are used, the more stable this sniper should be (theoretically)

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
