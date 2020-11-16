package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/shohinsherov/http/cmd/app"
	"github.com/shohinsherov/http/pkg/banners"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

type handler struct {
	mu       *sync.RWMutex
	handlers map[string]http.HandlerFunc
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.mu.RLock()
	handler, ok := h.handlers[request.URL.Path]
	h.mu.RUnlock()

	if !ok {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	handler(writer, request)
	/*_, err := writer.Write([]byte("Hello Bekhai be k"))
	if err != nil {
		log.Print(err)
	}*/
}

func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	
	
	server := app.NewServer(mux, bannersSvc)
	server.Init()

	/*mux.HandleFunc("/banners.getAll", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("Hello Bekhai be k"))
		if err != nil {
			log.Print(err)
		}
	}) */

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}

	log.Print("server start" + host + ":" + port)
	return srv.ListenAndServe()
}
