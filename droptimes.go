package main

import (
	"encoding/json"
	"errors"
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

type starShoppingResponseStruct struct {
	Unix int64 `json:"unix"`
}

func starShoppingDroptime(username string) (time.Time, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.star.shopping/droptime/%v", username), nil)

	if err != nil {
		return time.Time{}, err
	}

	req.Header.Set("User-Agent", "Sniper")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return time.Time{}, err
	}

	if resp.StatusCode < 300 {
		var starShoppingResponse starShoppingResponseStruct
		err = json.Unmarshal(respBytes, &starShoppingResponse)
		if err != nil {
			return time.Time{}, err
		}

		return time.Unix(starShoppingResponse.Unix, 0), nil
	}

	return time.Time{}, fmt.Errorf("failed to grab droptime with status %v and body \"%v\"", resp.Status, string(respBytes))
}

func getDroptime(username, preference string) (time.Time, error) {
	apis := map[string]func(string) (time.Time, error){
		"ckm":               coolkidmachoDroptime,
		"coolkidmacho":      coolkidmachoDroptime,
		"star.shopping":     starShoppingDroptime,
		"api.star.shopping": starShoppingDroptime,
	}
	allApis := []func(string) (time.Time, error){coolkidmachoDroptime, starShoppingDroptime}
	apisToUse := []func(string) (time.Time, error){}
	if val, ok := apis[preference]; ok {
		apisToUse = append(apisToUse, val)
	}

	apisToUse = append(apisToUse, allApis...)

	for _, api := range apisToUse {
		droptime, err := api(username)
		if err != nil {
			logErr("failed to grab droptime: %v", err)
			logInfo("trying next API")
			time.Sleep(time.Second * 1)
			continue
		}
		return droptime, nil
	}

	return time.Time{}, errors.New("failed to grab droptime from all APIs")
}

type next3RespStruct struct {
	Name string `json:"name"`
}
func getNext3c() ([]next3RespStruct, error){
	resp, err := http.Get("http://api.coolkidmacho.com/three")

	if err != nil {
		return []next3RespStruct{}, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []next3RespStruct{}, err
	}
	
	if resp.StatusCode < 300 {
		var respSlice []next3RespStruct
		err = json.Unmarshal(respBytes, &respSlice)
		if err != nil {
			return []next3RespStruct{}, err
		}
		return respSlice, nil
	}
	return []next3RespStruct{}, fmt.Errorf("failed to grab next 3c with status %v and body \"%v\"", resp.Status, string(respBytes))

}
