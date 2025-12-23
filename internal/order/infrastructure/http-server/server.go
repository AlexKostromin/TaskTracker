package http_server

import (
	"context"
	"net"
	"net/http"
	"time"

	modelOrder "gitlab.com/godevs2/micro/internal/order/domain/model"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router   *chi.Mux
	httpPort string
	server   *http.Server
	orderV1.UnimplementedHandler
	storage OrderProcessor
}
type OrderProcessor interface {
	CancelOrder(ctx context.Context, params modelOrder.CancelOrderParams) (orderV1.CancelOrderRes, error)
	CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error)
	GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error)
	PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error)
	NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode
}

/*func NewOrderHandler(storage OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}*/

func NewServer(httpPort string, orderHandler OrderProcessor) *Server {

	/*storage := mock-postgres.NewOrderStorage()
	orderService := application.NewOrderService(storage)
	orderHandler := NewOrderHandler(orderService)*/
	orderServer, err := orderV1.NewServer(orderHandler)
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
