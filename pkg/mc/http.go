package mc

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func (account *MCaccount) DefaultFastHttpHandler() {
	account.FastHttpClient = &fasthttp.Client{
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
		NoDefaultUserAgentHeader: true,
	}
}

func (account *MCaccount) SetProxy(proxy string) {

	if strings.HasPrefix(proxy, "socks") {
		account.FastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
	} else {
		if strings.HasPrefix(proxy, "http") {
			proxy = strings.TrimPrefix(proxy, "http://")
			proxy = strings.TrimPrefix(proxy, "https://")
		}
		account.FastHttpClient.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
	}

}
