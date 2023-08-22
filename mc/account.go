package mc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

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

	req, err := account.AuthenticatedReq("GET", url, nil)
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

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New(resp.Status)
}

func (account *MCaccount) AuthenticatedReq(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if account.Bearer == "" {
		return nil, errors.New("account is not authenticated")
	}
	req.Header.Add("Authorization", "Bearer "+account.Bearer)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (account *MCaccount) authenticate() error {
	payload := fmt.Sprintf(`{
    "agent": {                              
        "name": "Minecraft",                
        "version": 1                        
    },
    "username": "%s",      
    "password": "%s",
	"requestUser": true
}`, account.Email, account.Password)

	u := bytes.NewReader([]byte(payload))
	request, err := http.NewRequest("POST", "https://authserver.mojang.com/authenticate", u)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 300 {
		var AccountInfo authenticateReqResp
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &AccountInfo)
		if err != nil {
			return err
		}

		account.Bearer = AccountInfo.Accesstoken
		account.Username = AccountInfo.User.Username
		account.UUID = AccountInfo.User.ID
		return nil

	} else if resp.StatusCode == 403 {
		return errors.New("invalid email or password")
	}
	return errors.New("reached end of authenticate function! Shouldn't be possible. most likely 'failed to auth' status code changed")
}

func (account *MCaccount) loadSecurityQuestions() error {
	req, err := account.AuthenticatedReq("GET", "https://api.mojang.com/user/security/challenges", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("got status %v when requesting security questions", resp.Status)
	}

	defer resp.Body.Close()

	var sqAnswers []SqAnswer

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, &sqAnswers)
	if err != nil {
		return err
	}

	account.SecurityQuestions = sqAnswers

	return nil
}

// load account information (username, uuid) into accounts attributes, if not already there. When using Mojang authentication it is not necessary to load this info, as it will be automatically loaded.
func (account *MCaccount) LoadAccountInfo() error {
	req, err := account.AuthenticatedReq("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return errors.New("account does not own minecraft")
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var respJson accInfoResponse

	json.Unmarshal(respBytes, &respJson)

	account.Username = respJson.Name
	account.UUID = respJson.ID

	return nil
}

func (account *MCaccount) needToAnswer() (bool, error) {
	req, err := account.AuthenticatedReq("GET", "https://api.mojang.com/user/security/location", nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return false, nil
	}
	if resp.StatusCode == 403 {
		return true, nil
	}
	return true, fmt.Errorf("status of %v in needToAnswer not expected", resp.Status)
}

func (account *MCaccount) submitAnswers() error {
	if len(account.SecurityAnswers) != 3 {
		return errors.New("not enough security question answers provided")
	}
	if len(account.SecurityQuestions) != 3 {
		return errors.New("security questions not properly loaded")
	}
	var jsonContent []submitPostJson
	for i, sq := range account.SecurityQuestions {
		jsonContent = append(jsonContent, submitPostJson{ID: sq.Answer.ID, Answer: account.SecurityAnswers[i]})
	}
	jsonStr, err := json.Marshal(jsonContent)
	if err != nil {
		return err
	}
	req, err := account.AuthenticatedReq("POST", "https://api.mojang.com/user/security/location", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode == 204 {
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return errors.New("at least one security question answer was incorrect")
	}
	return fmt.Errorf("got status %v on post request for sqs", resp.Status)
}

// Runs all steps necessary to have a fully authenticated mojang account. It will submit email & pass and securitty questions (if necessary).
func (account *MCaccount) MojangAuthenticate() error {
	err := account.authenticate()
	if err != nil {
		return err
	}

	account.loadSecurityQuestions()

	if len(account.SecurityQuestions) == 0 {
		account.Authenticated = true
		return nil
	}

	answerNeeded, err := account.needToAnswer()
	if err != nil {
		return err
	}

	if !answerNeeded {
		account.Authenticated = true
		return nil
	}

	err = account.submitAnswers()
	if err != nil {
		return err
	}

	account.Authenticated = true
	return nil
}

func (account *MCaccount) HasGcApplied() (bool, error) {
	bodyStr := `{"profileName": "test"}`
	req, err := account.AuthenticatedReq("POST", "https://api.minecraftservices.com/minecraft/profile", bytes.NewReader([]byte(bodyStr)))
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return false, errors.New("successfully created profile with name test. unintended behavior, function is meant to check if gc is applied")

	} else if resp.StatusCode == 401 {
		return false, errors.New("received unauthorized response")
	} else if resp.StatusCode == 400 {
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

	return false, fmt.Errorf("got status: %v body: %v", resp.Status, string(bodyBytes))

}

// grab information on the availability of name change for this account
func (account *MCaccount) NameChangeInfo() (nameChangeInfoResponse, error) {
	req, err := account.AuthenticatedReq("GET", "https://api.minecraftservices.com/minecraft/profile/namechange", nil)

	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nameChangeInfoResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nameChangeInfoResponse{}, err
	}

	if resp.StatusCode >= 400 {
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

func (account *MCaccount) ChangeName(username string, changeTime time.Time, createProfile bool, proxy string) (NameChangeReturn, error) {
	client := &fasthttp.Client{
		Dial: fasthttp.Dial,
	}

	if proxy != "" {
		proxy = strings.TrimPrefix(proxy, "http://")
		proxy = strings.TrimPrefix(proxy, "https://")
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
	}

	var err error

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	if createProfile {
		req.Header.SetRequestURI("https://api.minecraftservices.com/minecraft/profile")
		req.Header.SetMethod("POST")
		req.SetBodyString(fmt.Sprintf(`{"profileName": "%s"}`, username))
	} else {
		req.Header.SetRequestURI(fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/name/%s", username))
		req.Header.SetMethod("PUT")
	}

	req.Header.Set("Authorization", "Bearer "+account.Bearer)
	req.Header.SetContentType("application/json")

	time.Sleep(time.Until(changeTime))

	sendTime := time.Now()
	err = client.Do(req, resp)
	recvTime := time.Now()

	if err != nil {
		return NameChangeReturn{Username: username}, err
	}

	return NameChangeReturn{
		Username:    username,
		ChangedName: resp.StatusCode() < 300,
		StatusCode:  resp.StatusCode(),
		SendTime:    sendTime,
		ReceiveTime: recvTime,
		Account:     *account,
	}, nil
}

func (account *MCaccount) ChangeSkinFromUrl(url, variant string) error {
	body := fmt.Sprintf(`{"url": "%v", "variant": "%v"}`, url, variant)
	req, err := account.AuthenticatedReq("POST", "https://api.minecraftservices.com/minecraft/profile/skins", strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("unauthorized")
	}

	return nil
}
