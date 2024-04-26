package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	backendManager "github.com/Kqzz/MCsniperGO/pkg/backend-manager"
	"github.com/Kqzz/MCsniperGO/pkg/claimer"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx            context.Context
	AccountManager *backendManager.AccountManager
	ProxyManager   *backendManager.ProxyManager
	QueueManager   *backendManager.QueueManager
	Claimer        *claimer.Claimer
}

// NewApp creates a new App application struct
func NewApp() *App {

	db, err := gorm.Open(sqlite.Open("mcsnipergo.db"), &gorm.Config{})
	db.AutoMigrate(&backendManager.Account{})
	db.AutoMigrate(&backendManager.Proxy{})
	db.AutoMigrate(&backendManager.Queue{})

	if err != nil {
		panic("failed to connect database")
	}

	claimer := &claimer.Claimer{}

	accountManager := backendManager.NewAccountManager()
	accountManager.DB = db
	accountManager.Claimer = claimer

	proxyManager := backendManager.NewProxyManager()
	proxyManager.DB = db
	proxyManager.Claimer = claimer

	queueManager := backendManager.NewQueueManager()
	queueManager.DB = db
	queueManager.Claimer = claimer

	return &App{
		AccountManager: accountManager,
		ProxyManager:   proxyManager,
		QueueManager:   queueManager,
		Claimer:        claimer,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbAccounts := a.AccountManager.GetAccounts()

	mcAccounts := []*mc.MCaccount{}

	for _, account := range dbAccounts {
		mcAcconunt := &mc.MCaccount{}
		Recast(account, mcAcconunt)
		mcAccounts = append(mcAccounts, mcAcconunt)
	}

	dbProxies := a.ProxyManager.GetProxies()

	// TODO: implement this better
	strProxies := []string{}

	for _, proxy := range dbProxies {
		prefix := ""
		if (proxy.Type == backendManager.SOCKS5 || proxy.Type == backendManager.SOCKS4) && !strings.HasPrefix(proxy.Url, "socks5://") {
			prefix = "socks5://" // TODO: figure out if fasthttpproxy supports SOCKS4 or not?
		}

		strProxies = append(strProxies, prefix+proxy.Url)
	}

	dialers := claimer.GetDialers(strProxies)

	a.Claimer.Accounts = mcAccounts
	a.Claimer.Dialers = dialers

	a.Claimer.Setup()

	dbQueues, _ := a.QueueManager.GetQueues()

	for _, queue := range dbQueues {
		a.Claimer.Queue(queue.Username, mc.DropRange{Start: time.Unix(queue.StartTime, 0), End: time.Unix(queue.EndTime, 0)})
	}
}

func Recast(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}
