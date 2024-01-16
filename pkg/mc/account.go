package mc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type FailType string

var (
	NOT_ENTITLED         FailType = "NOT_ENTITLED"
	DUPLICATE            FailType = "DUPLICATE"
	NOT_ALLOWED          FailType = "NOT_ALLOWED"
	CONSTRAINT_VIOLATION FailType = "CONSTRAINT_VIOLATION"
	TOO_MANY_REQUESTS    FailType = "TOO_MANY_REQUESTS"
)

func (account *MCaccount) AuthenticatedReq(method string, url string, body io.Reader) (*fasthttp.Request, *fasthttp.Response, error) {
	if account.Bearer == "" {
		return nil, nil, errors.New("bearer token not detected on account")
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.Header.SetRequestURI(url)
	req.Header.SetMethod(method)

	req.Header.Set("Authorization", "Bearer "+account.Bearer)
	req.Header.Set("Content-Type", "application/json")

	if body != nil {
		req.SetBodyStream(body, -1)
	}

	return req, resp, nil
}

// load account information (username, uuid) into accounts attributes, if not already there. When using Mojang authentication it is not necessary to load this info, as it will be automatically loaded.
func (account *MCaccount) LoadAccountInfo() error {
	req, resp, err := account.AuthenticatedReq("GET", "https://api.minecraftservices.com/minecraft/profile", nil)

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return err
	}

	err = account.FastHttpClient.Do(req, resp)

	if err != nil {
		return err
	}

	statusCode := resp.StatusCode()

	if statusCode == 404 {
		return errors.New("account does not own minecraft")
	}

	respBytes := resp.Body()

	if err != nil {
		return err
	}

	var respJson accInfoResponse

	json.Unmarshal(respBytes, &respJson)

	account.Username = respJson.Name
	account.UUID = respJson.ID

	return nil
}

func (account *MCaccount) HasGcApplied() (bool, error) {
	bodyStr := `{"profileName": "test"}`

	req, resp, err := account.AuthenticatedReq("POST", "https://api.minecraftservices.com/minecraft/profile", bytes.NewReader([]byte(bodyStr)))
	if err != nil {
		return false, err
	}

	err = account.FastHttpClient.Do(req, resp)
	if err != nil {
		return false, err
	}

	bodyBytes := resp.Body()
	if err != nil {
		return false, err
	}

	statusCode := resp.StatusCode()

	if statusCode == 200 {
		return false, errors.New("successfully created profile with name test. unintended behavior, function is meant to check if gc is applied")

	} else if statusCode == 401 {
		return false, errors.New("received unauthorized response")
	} else if statusCode == 400 {
		var respError hasGcAppliedResp

		err = json.Unmarshal(bodyBytes, &respError)
		if err != nil {
			return false, err
		}

		var hasGc bool

		switch respError.Details.Status {
		case "ALREADY_REGISTERED", "NOT_ENTITLED":
			{
				hasGc = false
			}
		case "DUPLICATE", "NOT_ALLOWED":
			{
				hasGc = true
			}
		default:
			{
				hasGc = false
			}
		}

		return hasGc, nil

	}

	if strings.Contains(string(bodyBytes), "Request blocked") {
		return false, errors.New("blocked by cloudfront (ip block)")
	}

	return false, fmt.Errorf("got status: %v body: %v", statusCode, string(bodyBytes))

}

// grab information on the availability of name change for this account
func (account *MCaccount) NameChangeInfo() (nameChangeInfoResponse, error) {
	req, resp, err := account.AuthenticatedReq("GET", "https://api.minecraftservices.com/minecraft/profile/namechange", nil)

	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	err = account.FastHttpClient.Do(req, resp)
	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	respBody := resp.Body()

	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	statusCode := resp.StatusCode()

	if statusCode >= 400 {
		return nameChangeInfoResponse{
			Changedat:         time.Time{},
			Createdat:         time.Time{},
			Namechangeallowed: false,
		}, errors.New("failed to grab name change info")
	}

	var parsedNameChangeInfo nameChangeInfoResponse

	err = json.Unmarshal(respBody, &parsedNameChangeInfo)

	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	return parsedNameChangeInfo, nil
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (account *MCaccount) License() error {
	url := fmt.Sprintf("https://api.minecraftservices.com/entitlements/license?requestId=%v", randomString(10))

	req, resp, err := account.AuthenticatedReq("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("authority", "api.minecraftservices.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.6")
	req.Header.Add("origin", "https://www.minecraft.net")
	req.Header.Add("referer", "https://www.minecraft.net/")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-gpc", "1")

	err = account.FastHttpClient.Do(req, resp)

	if err != nil {
		return err
	}

	statusCode := resp.StatusCode()

	if statusCode == 200 {
		return nil
	}

	return fmt.Errorf("failed w/ status: %v", statusCode)
}

func (account *MCaccount) CreateProfile(username string, client *fasthttp.Client) (int, FailType, error) {
	body := fmt.Sprintf(`{"profileName": "%s"}`, username)
	req, resp, err := account.AuthenticatedReq("POST", "https://api.minecraftservices.com/minecraft/profile", strings.NewReader(body))
	if err != nil {
		return 0, "", err
	}

	err = client.Do(req, resp)

	if err != nil {
		return 0, "", err
	}

	statusCode := resp.StatusCode()

	responseBody := resp.Body()

	if statusCode == 200 {
		return statusCode, "", nil
	}

	var fail FailType

	if statusCode == 429 {
		fail = TOO_MANY_REQUESTS
		return statusCode, fail, nil
	}

	for _, failType := range []FailType{NOT_ENTITLED, DUPLICATE, NOT_ALLOWED, CONSTRAINT_VIOLATION} {
		if strings.Contains(string(responseBody), string(failType)) {
			fail = failType
			break
		}
	}

	return statusCode, fail, nil
}
func (account *MCaccount) ChangeUsername(username string, client *fasthttp.Client) (int, FailType, error) {
	req, resp, err := account.AuthenticatedReq("PUT", fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/name/%s", username), nil)

	if err != nil {
		return 0, "", err
	}

	err = client.Do(req, resp)

	if err != nil {
		return 0, "", err
	}

	statusCode := resp.StatusCode()

	var fail FailType
	if statusCode == 200 {
		return statusCode, "", nil
	}

	if statusCode == 429 {
		fail = TOO_MANY_REQUESTS
		return statusCode, fail, nil
	}

	if statusCode == 403 {
		fail = DUPLICATE
	}

	return statusCode, fail, nil
}

func (account *MCaccount) ChangeSkinFromUrl(url, variant string) error {
	body := fmt.Sprintf(`{"url": "%v", "variant": "%v"}`, url, variant)
	req, resp, err := account.AuthenticatedReq("POST", "https://api.minecraftservices.com/minecraft/profile/skins", strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	err = account.FastHttpClient.Do(req, resp)

	if err != nil {
		return err
	}

	statusCode := resp.StatusCode()

	if statusCode != 200 {
		return fmt.Errorf("failed with status: %v", statusCode)
	}

	return nil
}
