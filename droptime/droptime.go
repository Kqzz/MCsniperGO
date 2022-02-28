package droptime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type droptimeFunction func(string) (time.Time, error)

type coolkidmachoRespStruct struct {
	Unix int64 `json:"UNIX"`
}

type starShoppingResponseStruct struct {
	Unix int64 `json:"unix"`
}

type threeChar struct {
	Name string `json:"name"`
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

// grabs droptime from api.star.shopping
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

	return time.Time{}, fmt.Errorf("failed to grab droptime with status %v and body %v", resp.Status, string(respBytes))
}

func GetDroptime(username string) (time.Time, error) {
	allApis := []droptimeFunction{starShoppingDroptime, coolkidmachoDroptime}

	errs := []string{}

	for _, api := range allApis {
		droptime, err := api(username)
		if err != nil {
			time.Sleep(time.Second * 1)
			errs = append(errs, err.Error())
			continue
		}
		return droptime, nil
	}

	return time.Time{}, errors.New(strings.Join(errs, "\n"))
}

func GetNext3c() ([]threeChar, error) {
	resp, err := http.Get("http://api.coolkidmacho.com/three")

	if err != nil {
		return []threeChar{}, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []threeChar{}, err
	}

	if resp.StatusCode < 300 {
		var respSlice []threeChar
		err = json.Unmarshal(respBytes, &respSlice)
		if err != nil {
			return []threeChar{}, err
		}
		return respSlice, nil
	}
	return []threeChar{}, fmt.Errorf("failed to grab next 3c with status %v and body %v", resp.Status, string(respBytes))
}
