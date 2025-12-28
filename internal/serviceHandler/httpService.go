package servicehandler

import (
	"fmt"
	"net/http"
	"scheduler/internal/config"
)

type HttpService struct {
	opts      *config.Options
	serverMux *http.ServeMux
}

func NewHttpService(opts *config.Options) *HttpService {
	httpService := HttpService{
		opts: opts,
	}

	httpService.initServer()
	return &httpService
}

func (h *HttpService) initServer() {
	serverMux := http.NewServeMux()

	// register handler function for each part
	serverMux.Handle("/", http.HandlerFunc(landing))
	serverMux.Handle("/ping", http.HandlerFunc(ping))

	// auth mux for use middleware for authentication
	authMux := http.NewServeMux()
	authMux.Handle("/ping", http.HandlerFunc(pingWithAuth))

	// strip prefix of /auth to be handler of auth middleware
	authHandler := http.StripPrefix("/auth", authMux)

	// set middleware
	serverMux.Handle("/auth/", authHandler)

	h.serverMux = serverMux
}

func (h *HttpService) Run() {
	addr := h.opts.GetAddress()
	fmt.Println("Run serve", addr)
	if err := http.ListenAndServe(addr, h.serverMux); err != nil {
		fmt.Println("Failed to run: ", err)
	}
}

func (h *HttpService) RegisterRoute() {
	// h.initServer()
}

func (h *HttpService) OnShutdown() {

}

// ////////////////////////////
// middlerware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth middleware")

		// TODO: handle of user

		next.ServeHTTP(w, r)
	})
}

// ////////////////////////////
// handler func
func landing(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Not found"))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func pingWithAuth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong with auth"))
}
