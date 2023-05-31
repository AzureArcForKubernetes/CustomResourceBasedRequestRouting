package requesthandler

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/AzureArcForKubernetes/CustomResourceBasedRequestRouting/controllers"
	"k8s.io/klog"
)

func StartProxyServer() {
	http.HandleFunc("/", redirectRequestBasedOnRules)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fatalError := fmt.Sprintf("Failed to start server: %s", err)
		klog.Fatalf(fatalError)
	}
}

func redirectRequestBasedOnRules(w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("%s%s", "https://test-request-router.rp.kubernetesconfiguration-test.azure.com", req.RequestURI)
	proxyReq, _ := http.NewRequest(req.Method, url, req.Body)
	//fmt.Printf("Request URL: %v\n", req)

	proxyReq.Header = req.Header
	//fmt.Printf("proxyReq URL: %v\n", proxyReq)
	httpClient := customClient(controllers.ResolveProxyEndpoint(req.RequestURI))
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("respBody URL: %v\n", string(respBody))

	w.Write(respBody)
	w.WriteHeader(resp.StatusCode)
}

func customClient(ip string) *http.Client {
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	//ref: Copy and modify defaults from https://golang.org/src/net/http/transport.go
	//Note: Clients and Transports should only be created once and reused
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  dialer.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "test-request-router.rp.kubernetesconfiguration-test.azure.com:443" {
				addr = ip + ":443"
			}
			return dialer.DialContext(ctx, network, addr)
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   4 * time.Second,
	}

	return &client
}
