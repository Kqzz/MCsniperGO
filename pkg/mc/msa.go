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
	"net/url"
	"regexp"
	"strings"
	"time"
)

type xBLSignInBody struct {
	Properties struct {
		Authmethod string `json:"AuthMethod"`
		Sitename   string `json:"SiteName"`
		Rpsticket  string `json:"RpsTicket"`
	} `json:"Properties"`
	Relyingparty string `json:"RelyingParty"`
	Tokentype    string `json:"TokenType"`
}

type XBLSignInResp struct {
	Issueinstant  time.Time `json:"IssueInstant"`
	Notafter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	Displayclaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type xSTSPostBody struct {
	Properties struct {
		Sandboxid  string   `json:"SandboxId"`
		Usertokens []string `json:"UserTokens"`
	} `json:"Properties"`
	Relyingparty string `json:"RelyingParty"`
	Tokentype    string `json:"TokenType"`
}

type xSTSAuthorizeResponse struct {
	Issueinstant  time.Time `json:"IssueInstant"`
	Notafter      time.Time `json:"NotAfter"`
	Token         string    `json:"Token"`
	Displayclaims struct {
		Xui []struct {
			Uhs string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

type xSTSAuthorizeResponseFail struct {
	Identity string `json:"Identity"`
	Xerr     int64  `json:"XErr"`
	Message  string `json:"Message"`
	Redirect string `json:"Redirect"`
}

type msGetMojangbearerBody struct {
	Identitytoken       string `json:"identityToken"`
	Ensurelegacyenabled bool   `json:"ensureLegacyEnabled"`
}

type msGetMojangBearerResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Foci         string `json:"foci"`
}

func (account *MCaccount) MicrosoftAuthenticate(proxy string) error {

	if account.Password == "code" {
		return account.OauthFlow()
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}

	if proxy != "" {
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	var redirect string
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirect = req.URL.String()
			return nil
		},
		Jar:       jar,
		Transport: tr,
	}
	// Grab value and urlpost
	valRegex := regexp.MustCompile(`value=\\\"(.+?)\\\"`)
	urlPostRegex := regexp.MustCompile(`urlPost":"(.+?)""`)

	resp, err := client.Get("https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en")

	if err != nil {
		return err
	}

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	value := string(valRegex.FindAllSubmatch(respBytes, -1)[0][1])
	urlPost := string(urlPostRegex.FindAllSubmatch(respBytes, -1)[0][1])

	// Sign in to microsoft

	emailEncoded := url.QueryEscape(account.Email)
	passwordEncoded := url.QueryEscape(account.Password)
	valueEncoded := url.QueryEscape(value)

	body := []byte(fmt.Sprintf("login=%v&loginfmt=%v&passwd=%v&PPFT=%v", emailEncoded, emailEncoded, passwordEncoded, valueEncoded))

	req, err := http.NewRequest("POST", urlPost, bytes.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.Request.URL.String() == urlPost && strings.Contains(resp.Request.URL.String(), "access_token") {
		return errors.New("invalid credentials, no access_token")
	}

	respBytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	respStr := string(respBytes)

	if strings.Contains(respStr, "Sign in to") {
		return errors.New("invalid credentials, sign in to")
	}

	if strings.Contains(respStr, "Help us protect your account") {
		return errors.New("2fa is enabled, which is not supported now")
	}

	if !strings.Contains(redirect, "access_token") || redirect == urlPost {
		return errors.New("invalid credentials, no access_token in redirect")
	}

	params := strings.Split(redirect, "#")[1]

	loginData := map[string]string{}

	for _, item := range strings.Split(params, "&") {
		itemSplit := strings.Split(item, "=")
		v, _ := url.QueryUnescape(itemSplit[1])
		loginData[itemSplit[0]] = v
	}

	data := xBLSignInBody{
		Properties: struct {
			Authmethod string "json:\"AuthMethod\""
			Sitename   string "json:\"SiteName\""
			Rpsticket  string "json:\"RpsTicket\""
		}{
			Authmethod: "RPS",
			Sitename:   "user.auth.xboxlive.com",
			Rpsticket:  loginData["access_token"],
		},
		Relyingparty: "http://auth.xboxlive.com",
		Tokentype:    "JWT",
	}

	encodedBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err = http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewReader(encodedBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-xbl-contract-version", "1")

	resp, err = client.Do(req)
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("login_with_xbox failed: %v", resp.Status)
	}

	var mcBearerResp msGetMojangBearerResponse

	json.Unmarshal(mcBearerResponseBytes, &mcBearerResp)

	account.Bearer = mcBearerResp.AccessToken

	return nil
}
