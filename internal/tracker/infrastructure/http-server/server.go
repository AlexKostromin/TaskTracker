package http_server

import (
	"context"
	"net"
	"net/http"
	"time"

	"log"

	models "github.com/AlexKostromin/TaskTracker/internal/tracker/domain"
	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router   *chi.Mux
	httpPort string
	server   *http.Server
	storage  TrackerProcessor
}
type TrackerProcessor interface {
	CreateTracker(ctx context.Context, request models.CreateTrackerRequest) (models.Tracker, error)
	UpdateTracker(ctx context.Context, request models.UpdateTrackerRequest, params models.UpdateTrackerParams) (models.Tracker, error)
	//NewError(ctx context.Context, err error) error
}

func NewServer(httpPort string, trackerHandler TrackerProcessor) *Server {

	s := &Server{
		httpPort: httpPort,
		storage:  trackerHandler,
	}

	trackerServer, err := trackerV1.NewServer(s)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Регистрация маршрутов
	r.Mount("/", trackerServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s.router = r

	return s
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:              net.JoinHostPort("localhost", s.httpPort),
		Handler:           s.router,
		ReadHeaderTimeout: 30 * time.Second, // Защита от Slowloris атак
	}

	log.Printf("🚀 HTTP-сервер запущен на порту %s\n", s.httpPort)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		log.Println("🛑 Остановка HTTP сервера...")
		return s.server.Shutdown(ctx)
	}
	return nil
}
