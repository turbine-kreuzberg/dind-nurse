package nurse

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

const MaxExecutionTime = 10 * time.Minute

// Service runs the wedding server.
type Service struct {
	router http.Handler
}

// NewService creates a new service server and initiates the routes.
func NewService(targetURL *url.URL) *Service {
	srv := &Service{}

	srv.routes(targetURL)

	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) routes(targetURL *url.URL) {
	router := mux.NewRouter()
	router.HandleFunc("/_nurse_healthy", ping).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(NewForwarder(targetURL))

	s.router = router
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewForwarder(targetURL *url.URL) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		/*
			// cpu limit
			cpuquota, err := strconv.Atoi(r.URL.Query().Get("cpuquota"))
			if err != nil {
				return cfg, fmt.Errorf("parse cpu quota to int: %v", err)
			}
			if cpuquota == 0 {
				cpuquota = buildCPUQuota
			}

			cpuperiod, err := strconv.Atoi(r.URL.Query().Get("cpuperiod"))
			if err != nil {
				return cfg, fmt.Errorf("parse cpu period to int: %v", err)
			}
			if cpuperiod == 0 {
				cpuperiod = buildCPUPeriod
			}
		*/

		proxy.ServeHTTP(w, r)
	}
}
