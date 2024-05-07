package proxy

import (
	"fmt"
	"github.com/gitlayzer/gee_loadbalancer/cmd/balancer"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	XServer       = http.CanonicalHeaderKey("X-Server")
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

// HTTPProxy 是一个 HTTP 代理
type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer

	sync.RWMutex
	alive map[string]bool
}

// NewHTTPProxy 新建一个 HTTP 代理
func NewHTTPProxy(targetNodes []string, algorithm string) (*HTTPProxy, error) {
	nodes := make([]string, 0) // nodes 是存放所有可用的节点的列表

	hostMap := make(map[string]*httputil.ReverseProxy) // hostMap 是存放所有节点的 map

	alive := make(map[string]bool) // alive 是存放节点是否存活的 map

	for _, targetNode := range targetNodes {
		u, err := url.Parse(targetNode)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(u)

		originDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
			req.Header.Set(XServer, "GeeLoadBalancer")
		}

		host := GetNode(u)
		alive[host] = true // initial mark alive
		hostMap[host] = proxy
		nodes = append(nodes, host)
	}

	loadBalance, err := balancer.Build(algorithm, nodes)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		hostMap: hostMap,
		lb:      loadBalance,
		alive:   alive,
	}, nil
}

// ServerHTTP 实现 http.Handler 接口
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic :%s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()

	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}

	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}
