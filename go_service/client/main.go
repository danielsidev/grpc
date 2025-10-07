// client/main.go
package main

import (
	"context"
	"log"
	"time"

	pb "grpc_comunication/simple" // Use o caminho do seu módulo

	"crypto/x509"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	// Adiciona o cabeçalho de autorização ao metadata
	return map[string]string{
		"authorization": t.token,
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	// Define se o TLS é necessário (deve ser true em produção)
	return true // Se estiver usando TLS
}

const (
	address = "localhost:50051"
)

func main() {
	// Seu token secreto (ex: um JWT)
	authToken := "Bearer my-secret-api-key-12345"

	// Configurar a autenticação de token
	auth := &tokenAuth{token: authToken}
	// Carregar o certificado público do servidor (CA raiz para o cliente)
	cert, err := os.ReadFile("server.pem")
	if err != nil {
		log.Fatalf("Falha ao ler o certificado CA: %v", err)
	}

	// Criar um pool de certificados com o CA do servidor
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		log.Fatalf("Falha ao adicionar certificado ao pool")
	}

	// Criar credenciais de transporte TLS
	tlsCreds := credentials.NewClientTLSFromCert(cp, "localhost")

	// 1. Conectar ao servidor gRPC
	// conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(tlsCreds), // Usando credenciais TLS do Exemplo 1
		grpc.WithPerRPCCredentials(auth),        // Adiciona o token à chamada RPC
	)
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	// 2. Criar o stub do cliente
	client := pb.NewOrderServiceClient(conn)

	// 3. Criar o payload de ENTRADA (equivalente ao body/json em REST)
	request := &pb.OrderRequest{
		ProductId: "GTX-3080",
		Quantity:  2,
		Price:     999.99,
	}

	// 4. Chamar o RPC (equivalente ao POST/GET em REST)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Enviando solicitação OrderRequest...")
	response, err := client.ProcessOrder(ctx, request)

	// 5. Tratar a resposta
	if err != nil {
		log.Fatalf("Falha na chamada RPC: %v", err)
	}

	// 6. Imprimir o payload de SAÍDA (o ReceiptResponse)
	log.Printf("--- Resposta (ReceiptResponse) ---")
	log.Printf("ID do Pedido: %s", response.GetOrderId())
	log.Printf("Valor Total:  %.2f", response.GetTotalAmount())
	log.Printf("Status:       %s", response.GetStatus())
	log.Printf("----------------------------------")
}
