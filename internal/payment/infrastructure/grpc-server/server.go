package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"gitlab.com/godevs2/micro/internal/payment/domain/model"
	paymentV1 "gitlab.com/godevs2/micro/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server   *grpc.Server
	grpcPort string
	paymentV1.UnimplementedPaymentServiceServer
	paymentService PaymentProcessor
}

type PaymentProcessor interface {
	Pay(ctx context.Context, req *model.PayOrder) (*model.PayOrderResponse, error)
}

func NewServer(grpcPort string, paymentHandler PaymentProcessor) *Server {
	// Создание gRPC сервера с интерсепторами
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggerInterceptor(),
		),
	)
	s := &Server{
		server:         server,
		grpcPort:       grpcPort,
		paymentService: paymentHandler,
	}
	//storage := mock.NewPaymentStorage() // S (solid)
	//paymentService := application.NewPaymentService(storage)
	//paymentHandler := NewPaymentHandler(paymentService)
	paymentV1.RegisterPaymentServiceServer(server, s)
	reflection.Register(server)

	return s
}

func (s *Server) Start() error {
	// Создаем listener для gRPC сервера
	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("🚀 gRPC server listening on %s\n", s.grpcPort) //zap /
	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	log.Println("🛑 Остановка gRPC сервера...")

	// Останавливаем gRPC сервер
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	// Ожидаем завершения остановки или таймаута
	select {
	case <-stopped:
		log.Println("✅ gRPC сервер остановлен")
	case <-ctx.Done():
		log.Println("⚠️  Таймаут graceful shutdown, принудительная остановка")
		s.server.Stop()
	}
	return nil
}

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Логируем входящий запрос
		log.Printf("➡️  gRPC call: %s", info.FullMethod)

		// Вызываем следующий обработчик
		resp, err := handler(ctx, req)

		// Логируем результат
		if err != nil {
			log.Printf("❌ gRPC errors: %s - %v", info.FullMethod, err)
		} else {
			log.Printf("✅ gRPC success: %s", info.FullMethod)
		}

		return resp, err
	}
}
