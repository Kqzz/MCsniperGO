package accountmanager

import (
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"gorm.io/gorm"
)

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

type AccountManager struct {
	DB *gorm.DB
}

type Account struct {
	gorm.Model
	Email    string
	Password string
	Username string
	Type     mc.AccType
	Bearer   string
}

func (am *AccountManager) AddAccount(account *Account) {
	am.DB.Create(account)
}

func (am *AccountManager) RemoveAccount(account *Account) {
	am.DB.Delete(account)
}

func (am *AccountManager) GetAccounts() []*Account {
	var accounts []*Account
	am.DB.Find(&accounts)
	return accounts
}
