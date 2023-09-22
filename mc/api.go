package mc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type ProfileResponse struct {
	Name string `json:"name"` // account username
	ID   string `json:"id"`   // UUID of account
}

func UsernameToUuid(username string) (ProfileResponse, int, error) {
	var profile ProfileResponse

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username))

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Fatalf("Request failed: %s", err)
		return profile, 0, err
	}

	statusCode := resp.StatusCode()

	bodyBytes := resp.Body()

	err = json.Unmarshal(bodyBytes, &profile)
	if err != nil {
		log.Fatalf("Unmarshalling failed: %s", err)
		return profile, statusCode, err
	}

	return profile, statusCode, nil
}
