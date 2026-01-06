package mc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

/*

Credit to emily (@impliedgg) for the entire oauth flow!
Client ID is 00000000441cc96b, Minecraft for Nintendo Switch

Flow is as follows:
POST https://login.live.com/oauth20_connect.srf
?client_id={client_id}
&scope=XboxLive.signin

Inform user to visit link from response.verification_uri and enter code response.user_code.

POST https://login.live.com/oauth20_token.srf
?grant_type=urn:ietf:params:oauth:grant-type:device_code
&client_id={client_id}
&device_code={respone.device_code}

once every response.interval seconds until expires_in timeout or successful poll.

Errors to properly handle in response.error:
authorization_pending - keep waiting. user isn't done.
authorization_declined - user declined auth, fail to authenticate.
bad_verification_code - this one should request a bug report on github. won't happen normally
expired_token - stop polling, fail to authenticate. user took too long.const

Fields to use once response.error is nil:
access_token - use this with https://user.auth.xboxlive.com/user/authenticate to get xsts done.
expires_in - if implemented, should request reauthentication once expired.

*/

// we only take the useful fields here.

type msDeviceInitResponse struct {
	VerificationURI string `json:"verification_uri"`
	UserCode        string `json:"user_code"`
	Message         string `json:"message"`
	Interval        int    `json:"interval"`
	DeviceCode      string `json:"device_code"`
}

type msErrorPollResponse struct {
	Error string `json:"error"`
}

type msSuccessPollResponse struct {
	AccessToken string `json:"access_token"`
}

const client_id = "00000000441cc96b"

// types in msa.go are used here as well.

func (account *MCaccount) OauthFlow() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: tr,
	}

	reqParams := fmt.Sprintf("client_id=%s&scope=XboxLive.signin&response_type=device_code", client_id)

	req, _ := http.NewRequest("POST", "https://login.live.com/oauth20_connect.srf", bytes.NewBuffer([]byte(reqParams)))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("non-200 status on devicecode post")
	}

	if err != nil {
		return err
	}

	var respObj msDeviceInitResponse
	err = json.Unmarshal(respbytes, &respObj)
	if err != nil {
		return err
	}
	fmt.Printf("[*] auth: Please visit %v and use the code %v to continue\n", respObj.VerificationURI, respObj.UserCode)

	return pollEndpoint(account, respObj.DeviceCode, respObj.Interval)
}

func authWithToken(account *MCaccount, access_token_from_ms string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: tr,
	}
	data := xBLSignInBody{
		Properties: struct {
			Authmethod string "json:\"AuthMethod\""
			Sitename   string "json:\"SiteName\""
			Rpsticket  string "json:\"RpsTicket\""
		}{
			Authmethod: "RPS",
			Sitename:   "user.auth.xboxlive.com",
			Rpsticket:  "d=" + access_token_from_ms,
		},
		Relyingparty: "http://auth.xboxlive.com",
		Tokentype:    "JWT",
	}

	encodedBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewReader(encodedBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-xbl-contract-version", "1")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		return errors.New("invalid Rpsticket field probably")
	}

	if err != nil {
		return err
	}

	var respBody XBLSignInResp

	json.Unmarshal(respBodyBytes, &respBody)

	uhs := respBody.Displayclaims.Xui[0].Uhs
	XBLToken := respBody.Token

	xstsBody := xSTSPostBody{
		Properties: struct {
			Sandboxid  string   "json:\"SandboxId\""
			Usertokens []string "json:\"UserTokens\""
		}{
			Sandboxid: "RETAIL",
			Usertokens: []string{
				XBLToken,
			},
		},
		Relyingparty: "rp://api.minecraftservices.com/",
		Tokentype:    "JWT",
	}

	encodedXstsBody, err := json.Marshal(xstsBody)
	if err != nil {
		return err
	}
	req, err = http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewReader(encodedXstsBody))
	if err != nil {
		return err
	}

	resp, err = client.Do(req)

	if err != nil {
		return err
	}

	respBodyBytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode == 401 {
		var authorizeXstsFail xSTSAuthorizeResponseFail
		json.Unmarshal(respBodyBytes, &authorizeXstsFail)
		switch authorizeXstsFail.Xerr {
		case 2148916238:
			{
				return errors.New("microsoft account belongs to someone under 18! add to family for this to work")
			}
		case 2148916233:
			{
				return errors.New("you have no xbox account! Sign up for one to continue")
			}
		default:
			{
				return fmt.Errorf("got error code %v when trying to authorize XSTS token", authorizeXstsFail.Xerr)
			}
		}
	}

	var xstsAuthorizeResp xSTSAuthorizeResponse
	json.Unmarshal(respBodyBytes, &xstsAuthorizeResp)

	xstsToken := xstsAuthorizeResp.Token

	mojangBearerBody := msGetMojangbearerBody{
		Identitytoken:       "XBL3.0 x=" + uhs + ";" + xstsToken,
		Ensurelegacyenabled: true,
	}

	mojangBearerBodyEncoded, err := json.Marshal(mojangBearerBody)

	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewReader(mojangBearerBodyEncoded))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	mcBearerResponseBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var mcBearerResp msGetMojangBearerResponse

	json.Unmarshal(mcBearerResponseBytes, &mcBearerResp)

	account.Bearer = mcBearerResp.AccessToken

	return nil
}

func pollEndpoint(account *MCaccount, device_code string, interval int) error {

	sleepDuration := time.Second * time.Duration(interval)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: tr,
	}

	reqParams := fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:device_code&device_code=%s&client_id=%s", device_code, client_id)
	for {
		time.Sleep(sleepDuration)
		req, err := http.NewRequest("POST", "https://login.live.com/oauth20_token.srf", bytes.NewBuffer([]byte(reqParams)))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		byteRes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode == 400 {
			var r msErrorPollResponse
			err = json.Unmarshal(byteRes, &r)
			if err != nil {
				return err
			}
			switch r.Error {
			case "authorization_pending":
				continue
			case "authorization_declined", "expired_token":
				return errors.New("authorization failed. cannot continue")
			default:
				return errors.New("unknown state on 400 status")
			}
		} else if resp.StatusCode == 200 {
			var r msSuccessPollResponse
			err = json.Unmarshal(byteRes, &r)
			if err != nil {
				return err
			}
			return authWithToken(account, r.AccessToken)
		} else {
			return errors.New("status code response not 200 or 400")
		}
	}
}
