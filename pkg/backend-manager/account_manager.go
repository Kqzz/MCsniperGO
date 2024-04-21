package backendmanager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Kqzz/MCsniperGO/pkg/claimer"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/Kqzz/MCsniperGO/pkg/parser"
	"gorm.io/gorm"
)

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

type AccountManager struct {
	DB      *gorm.DB
	Claimer *claimer.Claimer
}

type Account struct {
	gorm.Model
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Username string     `json:"username"`
	Type     mc.AccType `json:"type"`
	Bearer   string     `json:"bearer"`
}

func (am *AccountManager) AddAccounts(accounts string, accountType mc.AccType) error {

	lines := strings.Split(accounts, "\n")

	if len(lines) == 0 {
		return errors.New("no accounts to add")
	}

	mcAccounts, errs := parser.ParseAccounts(lines, accountType)

	if len(errs) > 0 {
		fmt.Println("Errors parsing accounts:", errs)
	}

	parsedAccounts := []*Account{}

	for _, acc := range mcAccounts {
		parsedAccounts = append(parsedAccounts, &Account{
			Email:    acc.Email,
			Password: acc.Password,
			Username: acc.Username,
			Type:     accountType,
		})
	}

	tx := am.DB.Create(parsedAccounts)
	fmt.Println("Created accounts:", parsedAccounts)

	return tx.Error

}

func (am *AccountManager) RemoveAccountByEmail(email string) {
	am.DB.Where("email = ?", email).Delete(&Account{})
}

func (am *AccountManager) GetAccounts() []*Account {
	var accounts []*Account
	am.DB.Find(&accounts)
	return accounts
}
