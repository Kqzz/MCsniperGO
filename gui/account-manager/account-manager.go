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
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Username string     `json:"username"`
	Type     mc.AccType `json:"type"`
	Bearer   string     `json:"bearer"`
}

func (am *AccountManager) AddAccount(account *Account) {
	am.DB.Create(account)
}

func (am *AccountManager) RemoveAccountByEmail(email string) {
	am.DB.Where("email = ?", email).Delete(&Account{})
}

func (am *AccountManager) GetAccounts() []*Account {
	var accounts []*Account
	am.DB.Find(&accounts)
	return accounts
}
