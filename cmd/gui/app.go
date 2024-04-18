package main

import (
	"context"

	backendManager "github.com/Kqzz/MCsniperGO/pkg/backend-manager"
	"github.com/Kqzz/MCsniperGO/pkg/claimer"
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

	accountManager := backendManager.NewAccountManager()
	accountManager.DB = db

	proxyManager := backendManager.NewProxyManager()
	proxyManager.DB = db

	queueManager := backendManager.NewQueueManager()
	queueManager.DB = db

	claimer := &claimer.Claimer{}

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
	a.Claimer.Setup()
}
