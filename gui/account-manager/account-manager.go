package accountmanager

import "github.com/Kqzz/MCsniperGO/pkg/mc"

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

type AccountManager struct {
	Accounts []*mc.MCaccount `json:"accounts"`
}
