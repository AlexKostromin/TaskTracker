package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"gitlab.com/godevs2/micro/internal/inventory/domain/model"
	inventoryV1 "gitlab.com/godevs2/micro/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server   *grpc.Server
	grpcPort string
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService InventoryProcessor
}

type InventoryProcessor interface {
	Get(ctx context.Context, req *model.GetPartRequest) (*model.GetPartResponse, error)
	ListParts(ctx context.Context, req *model.ListPartsRequest) (*model.ListPartsResponse, error)
}

func NewServer(grpcPort string, inventoryHandler InventoryProcessor) *Server {

	// Создание gRPC сервера с интерсепторами
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggerInterceptor(),
		),
	)
	// Регистрация сервисов

	s := &Server{
		server:           server,
		grpcPort:         grpcPort,
		inventoryService: inventoryHandler,
	}
	inventoryV1.RegisterInventoryServiceServer(server, s)
	reflection.Register(server)

	return s
}

func (s *Server) Start() error {
	// Создаем listener для gRPC сервера
	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("🚀 gRPC server listening on %s\n", s.grpcPort)
	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
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
