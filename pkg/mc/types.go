package mc

import (
	"time"

	"github.com/valyala/fasthttp"
)

type DropRange struct {
	Start time.Time
	End   time.Time
}

type AccType string

const (
	Ms   AccType = "MS"
	MsGc AccType = "GC"
)

/// RETURNS  / API-FACING ///
// Holds name change information for an account, the time the current account was created, it's name was most recently changed, and if it can currently change its name.

type NameChangeReturn struct {
	Email       string
	Account     MCaccount `json:"-"`
	Username    string
	ChangedName bool
	StatusCode  int
	SendTime    time.Time
	ReceiveTime time.Time
}

// represents a minecraft account
type MCaccount struct {
	Email          string           `json:"email"`
	Password       string           `json:"password"`
	Bearer         string           `json:"bearer"`
	BearerAccount  bool             `json:"bearerAccount"`
	RefreshToken   string           `json:"refreshToken"`
	UUID           string           `json:"uuid"`
	Username       string           `json:"username"`
	FastHttpClient *fasthttp.Client // client is used for all requests except create auth, profile create, and name change
	Type           AccType          `json:"type"`
}

/// HTTP RESPONSE BODIES ///

type nameChangeInfoResponse struct {
	Changedat         time.Time `json:"changedAt"`
	Createdat         time.Time `json:"createdAt"`
	Namechangeallowed bool      `json:"nameChangeAllowed"`
}

type hasGcAppliedResp struct {
	Path             string `json:"path"`
	ErrorType        string `json:"errorType"`
	Error            string `json:"error"`
	ErrorMessage     string `json:"errorMessage"`
	DeveloperMessage string `json:"developerMessage"`
	Details          struct {
		Status string `json:"status"`
	} `json:"details"`
}

type SqAnswer struct {
	Answer struct {
		ID int `json:"id"`
	} `json:"answer"`
	Question struct {
		ID       int    `json:"id"`
		Question string `json:"question"`
	} `json:"question"`
}

type accInfoResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
