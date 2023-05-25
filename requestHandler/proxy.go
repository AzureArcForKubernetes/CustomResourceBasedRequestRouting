package requesthandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AzureArcForKubernetes/CustomResourceBasedRequestRouting/controllers"
	"github.com/gorilla/mux"
	"k8s.io/klog"
)

func StartProxyServer() {
	muxRouter := getRouter()
	server := &http.Server{
		Addr:              ":8082",
		Handler:           muxRouter,
		ReadHeaderTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		fatalError := fmt.Sprintf("Failed to start server: %s", err)
		klog.Fatalf(fatalError)
	}
}

func getRouter() *mux.Router {
	serverMux := mux.NewRouter()
	s := serverMux.PathPrefix("/extensions").Subrouter()
	s.HandleFunc("/{extensionName}/{operationId}/events", http.HandlerFunc(redirectRequestBasedOnRules))
	return serverMux
}

func redirectRequestBasedOnRules(w http.ResponseWriter, req *http.Request) {
	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	httpClient := http.Client{}
	url := fmt.Sprintf("%s://%s%s", "https", controllers.ResolveProxyEndpoint(req.RequestURI), req.RequestURI)
	proxyReq, _ := http.NewRequest(req.Method, url, req.Body)
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
}
