package claimer

import (
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func GetDialers(proxies []string) []fasthttp.DialFunc {
	var dialers []fasthttp.DialFunc
	for _, proxy := range proxies {
		if strings.HasPrefix(proxy, "http") {
			dialer := fasthttpproxy.FasthttpHTTPDialer(proxy)
			dialers = append(dialers, dialer)
		} else {
			dialer := fasthttpproxy.FasthttpSocksDialer(proxy)
			dialers = append(dialers, dialer)
		}
	}

	if len(dialers) < 100 {
		dialers = append(dialers, fasthttp.Dial)
	}

	return dialers
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
