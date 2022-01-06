package nurse

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

const MaxExecutionTime = 60 * time.Minute

// Service runs the wedding server.
type Service struct {
	router http.Handler
}

// NewService creates a new service server and initiates the routes.
func NewService(targetURL *url.URL, dindMemoryLimit int) *Service {
	srv := &Service{}

	srv.routes(targetURL, dindMemoryLimit)

	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) routes(targetURL *url.URL, dindMemoryLimit int) {
	router := mux.NewRouter()
	router.HandleFunc("/_nurse_healthy", ping).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(NewForwarder(targetURL, dindMemoryLimit))

	s.router = router
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewForwarder(targetURL *url.URL, dindMemoryLimit int) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	bottleneck := &sync.Mutex{}
	openConnections := int64(0)

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

		bottleneck.Lock()
		atomic.AddInt64(&openConnections, 1)
		defer atomic.AddInt64(&openConnections, -1)
		bottleneck.Unlock()

		proxy.ServeHTTP(w, r)

		bottleneck.Lock()
		defer bottleneck.Unlock()
		if atomic.LoadInt64(&openConnections) == 1 {
			Cleanup(dindMemoryLimit)
		}
	}
}

func Cleanup(dindMemoryLimit int) {
	err := AvoidOOM(dindMemoryLimit)
	if err != nil {
		log.Printf("avoid oom: %v", err)
	}
}
