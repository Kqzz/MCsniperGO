package main

import (
	"embed"

	accountmanager "github.com/Kqzz/MCsniperGO/pkg/account-manager"
	"github.com/glebarez/sqlite"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gorm.io/gorm"
)

//go:embed all:frontend/dist
var assets embed.FS

const DEV bool = true

func main() {
	// Create an instance of the app structure
	app := NewApp()

	db, err := gorm.Open(sqlite.Open("mcsnipergo.db"), &gorm.Config{})
	db.AutoMigrate(&accountmanager.Account{})

	if err != nil {
		panic("failed to connect database")
	}

	accountManager := accountmanager.NewAccountManager()

	accountManager.DB = db

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "MCsniperGO",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		StartHidden:      DEV,
		Bind: []interface{}{
			app,
			accountManager,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
