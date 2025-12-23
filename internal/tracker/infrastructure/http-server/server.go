package http_server

import (
	"context"
	"net"
	"net/http"
	"time"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router    *chi.Mux
	httpPort  string
	server    *http.Server
	trackerv1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
	storage   TrackerProcessor
}
type TrackerProcessor interface {
}

/*func NewOrderHandler(storage OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}*/

func NewServer(httpPort string, trackerHandler TrackerProcessor) *Server {

	/*storage := mock-postgres.NewOrderStorage()
	orderService := application.NewOrderService(storage)
	orderHandler := NewOrderHandler(orderService)*/
	orderServer, err := trackerV1.NewServer(trackerHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Регистрация маршрутов
	r.Mount("/", orderServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	s := &Server{
		router:   r,
		httpPort: httpPort,
		/*orderService: orderHandler,*/
	}

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

/*func (s *Server) registerRoutes(r *chi.Mux) {
	// Регистрация всех обработчиков
	r.Get("/orders/{id}", s.orderStorage.GetOrder)
	r.Post("/orders", CreateOrderHandler)
	r.Post("/orders/{id}/pay", PayOrderHandler)
	r.Post("/orders/{id}/cancel", CancelOrderHandler)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
*/
