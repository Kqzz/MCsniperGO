package main

import (
	"errors"
	"strconv"
	"time"
)

// Feel free to add your droptime API to this file, the more droptime APIS that are used, the more stable this sniper should be (theoretically)

func manualDroptime(username string) (time.Time, error) {
	droptimeStr := userInput("enter unix droptime for %v", username)

	droptimeInt, err := strconv.ParseInt(droptimeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(droptimeInt, 0), nil
}

func getDroptime(username, preference string) (time.Time, error) {
	apis := map[string]func(string) (time.Time, error){
		"manual": manualDroptime,
	}
	allApis := []func(string) (time.Time, error){manualDroptime}
	apisToUse := []func(string) (time.Time, error){}
	if val, ok := apis[preference]; ok {
		apisToUse = append(apisToUse, val)
	}

	apisToUse = append(apisToUse, allApis...)

	for _, api := range apisToUse {
		droptime, err := api(username)
		if err != nil {
			log("error", "failed to grab droptime: %v", err)
			log("info", "trying next API")
			time.Sleep(time.Second * 1)
			continue
		}
		return droptime, nil
	}

	return time.Time{}, errors.New("failed to grab droptime from all APIs")
}
