package backendmanager

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/pkg/claimer"
	"gorm.io/gorm"
)

func NewProxyManager() *ProxyManager {
	return &ProxyManager{}
}

type ProxyManager struct {
	DB      *gorm.DB
	Claimer *claimer.Claimer
}

type ProxyType string

const (
	HTTP   ProxyType = "http"
	SOCKS4 ProxyType = "socks4"
	SOCKS5 ProxyType = "socks5"
)

type Proxy struct {
	gorm.Model
	Url  string    `json:"url"`
	Type ProxyType `json:"type"`
}

func (pm *ProxyManager) AddProxies(urls []string, proxyType ProxyType) error {
	fmt.Println("Adding proxies:", urls)
	proxies := []*Proxy{}

	for _, p := range urls {
		proxies = append(proxies, &Proxy{Url: p, Type: proxyType})
	}

	fmt.Println("Proxies:", proxies)

	tx := pm.DB.Create(proxies)
	return tx.Error
}

func (pm *ProxyManager) GetProxies() []Proxy {
	var proxies []Proxy
	pm.DB.Find(&proxies)

	return proxies
}

func (pm *ProxyManager) RemoveProxies(urls []string) error {
	proxies := []*Proxy{}

	for _, p := range urls {
		proxies = append(proxies, &Proxy{Url: p})
	}

	tx := pm.DB.Delete(&proxies, "url in ?", urls)
	return tx.Error
}
