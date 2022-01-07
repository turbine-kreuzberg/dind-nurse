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
	"golang.org/x/sync/semaphore"
)

const MaxExecutionTime = 60 * time.Minute

// Service runs the wedding server.
type Service struct {
	router http.Handler
}

// NewService creates a new service server and initiates the routes.
func NewService(targetURL *url.URL, dindMemoryLimit, parallelRequestLimit int, dockerPath string, diskUsageLimit int) *Service {
	srv := &Service{}

	srv.routes(targetURL, dindMemoryLimit, parallelRequestLimit, dockerPath, diskUsageLimit)

	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) routes(targetURL *url.URL, dindMemoryLimit, parallelRequestLimit int, dockerPath string, diskUsageLimit int) {
	router := mux.NewRouter()
	router.HandleFunc("/_nurse_healthy", ping).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(newForwarder(targetURL, dindMemoryLimit, parallelRequestLimit, dockerPath, diskUsageLimit))

	s.router = router
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func newForwarder(targetURL *url.URL, dindMemoryLimit, parallelRequestLimit int, dockerPath string, diskUsageLimit int) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	bottleneck := &sync.Mutex{}
	openConnections := int64(0)
	sem := semaphore.NewWeighted(int64(parallelRequestLimit * 2)) // times 2 to work with multiple connections triggered via DOCKER_BUILDKIT=1

	return func(w http.ResponseWriter, r *http.Request) {
		err := sem.Acquire(r.Context(), 1)
		if err != nil {
			log.Printf("failed to acquire semaphore slot: %v", err)
			return
		}
		defer sem.Release(1)

		bottleneck.Lock()
		atomic.AddInt64(&openConnections, 1)
		defer atomic.AddInt64(&openConnections, -1)
		bottleneck.Unlock()

		proxy.ServeHTTP(w, r)

		bottleneck.Lock()
		defer bottleneck.Unlock()
		if atomic.LoadInt64(&openConnections) == 1 {
			Cleanup(dindMemoryLimit, dockerPath, diskUsageLimit)
		}
	}
}

func Cleanup(dindMemoryLimit int, dockerPath string, diskUsageLimit int) {
	err := AvoidOOM(dindMemoryLimit)
	if err != nil {
		log.Printf("avoid oom: %v", err)
	}

	err = CollectGargabe(dockerPath, diskUsageLimit)
	if err != nil {
		log.Printf("avoid full disk: %v", err)
	}
}
