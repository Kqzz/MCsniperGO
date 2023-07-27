package mc

import "time"

type AccType string

const (
	Ms   AccType = "ms"
	Mj   AccType = "mj"
	MsPr AccType = "mspr"
)

/// RETURNS  / API-FACING ///

// Holds name change information for an account, the time the current account was created, it's name was most recently changed, and if it can currently change its name.

type NameChangeReturn struct {
	Account     MCaccount
	Username    string
	ChangedName bool
	StatusCode  int
	SendTime    time.Time
	ReceiveTime time.Time
}

// represents a minecraft account
type MCaccount struct {
	Email             string
	Password          string
	SecurityQuestions []SqAnswer
	SecurityAnswers   []string
	Bearer            string
	UUID              string
	Username          string
	Type              AccType
	Authenticated     bool
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

type authenticateReqResp struct {
	User struct {
		Properties []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"properties"`
		Username string `json:"username"`
		ID       string `json:"id"`
	} `json:"user"`
	Accesstoken string `json:"accessToken"`
	Clienttoken string `json:"clientToken"`
}

/// SEND BODIES ///
type submitPostJson struct {
	ID     int    `json:"id"`
	Answer string `json:"answer"`
}
