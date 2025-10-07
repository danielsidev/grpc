package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "grpc_comunication/simple" // Use o caminho do seu módulo

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

// orderServer é a implementação do nosso serviço OrderService.
type orderServer struct {
	pb.UnimplementedOrderServiceServer
}

// ProcessOrder implementa o RPC Unary ProcessOrder.
func (*orderServer) ProcessOrder(ctx context.Context, req *pb.OrderRequest) (*pb.ReceiptResponse, error) {
	log.Printf("Recebido pedido: ID: %s, Qtd: %d, Preço: %.2f",
		req.GetProductId(), req.GetQuantity(), req.GetPrice())

	// 1. Simulação de processamento de negócio
	total := float64(req.GetQuantity()) * req.GetPrice()
	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

	// 2. Criação do payload de resposta
	response := &pb.ReceiptResponse{
		OrderId:     orderID,
		TotalAmount: total,
		Status:      "SUCCESS",
	}

	return response, nil
}

var expectedToken string = os.Getenv("AUTH_TOKEN")

// Função de verificação do token
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := authorization[0]
	// Apenas um check simples: o token deve ser "Bearer my-secret-api-key-12345"
	log.Printf("Token recebido: %s", token)
	return token == expectedToken
}

// Interceptor para verificar a autenticação
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 1. Pegar o metadata do contexto
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata ausente")
	}

	// 2. Verificar o token de autorização
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "Token de autorização inválido ou ausente")
	}

	// 3. Se autenticado, chama o método RPC (ProcessOrder)
	return handler(ctx, req)
}

func main() {
	// Carregar credenciais TLS
	creds, err := credentials.NewServerTLSFromFile("server.pem", "server.key")
	if err != nil {
		log.Fatalf("Falha ao carregar credenciais TLS: %v", err)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Falha ao iniciar o listener: %v", err)
	}

	// Configurar o servidor gRPC com as credenciais
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authInterceptor))

	// Registrar o serviço no servidor gRPC
	pb.RegisterOrderServiceServer(s, &orderServer{})

	log.Printf("Servidor gRPC escutando em %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
