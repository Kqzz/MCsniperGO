package claimer

import (
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// TODO: reimplement this as an internal function to Claimer. The outside code should not be forced to manage the dial functions of the Claimer object.
func GetDialers(proxies []string) []fasthttp.DialFunc {
	var dialers []fasthttp.DialFunc
	for _, proxy := range proxies { // TODO: make it so proxies can create multiple threads per proxy (for rotating proxies), perhaps minimum-proxies configuration option?
		var dialer fasthttp.DialFunc
		if strings.HasPrefix(proxy, "socks5://") {
			dialer = fasthttpproxy.FasthttpSocksDialer(proxy)
		} else {
			dialer = fasthttpproxy.FasthttpHTTPDialer(proxy)
		}
		dialers = append(dialers, dialer)
	}

	if len(dialers) < 100 { // uses home IP if many proxies are not supplied
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
