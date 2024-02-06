package backendmanager

import (
	"gorm.io/gorm"
)

func NewProxyManager() *AccountManager {
	return &AccountManager{}
}

type ProxyManager struct {
	DB *gorm.DB
}

type Proxy struct {
	gorm.Model
	Url string `json:"url"`
}

func (pm *ProxyManager) AddProxies(urls []string) error {
	proxies := []*Proxy{}

	for _, p := range urls {
		proxies = append(proxies, &Proxy{Url: p})
	}

	tx := pm.DB.Create(&proxies)
	return tx.Error
}

func (pm *ProxyManager) GetProxies() ([]string, error) {
	var proxies []Proxy
	tx := pm.DB.Find(&proxies)
	return nil, tx.Error
}

func (pm *ProxyManager) DeleteProxies(urls []string) error {
	proxies := []*Proxy{}
	tx := pm.DB.Delete(&proxies, "url in ?", urls)
	return tx.Error
}
