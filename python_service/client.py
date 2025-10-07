import grpc
import os
import simple_pb2
import simple_pb2_grpc

# --- Configurações ---
SERVER_ADDRESS = 'localhost:50051'
CERT_PATH = 'server.pem'  # Certificado público do servidor Go
AUTH_TOKEN = 'Bearer my-secret-api-key-12345'  # Token que o servidor Go espera

def run():
    # 1. Carregar o certificado público (TLS)
    try:
        with open(CERT_PATH, 'rb') as f:
            root_certs = f.read()
    except FileNotFoundError:
        print(f"ERRO: Arquivo '{CERT_PATH}' não encontrado. Certifique-se de copiá-lo.")
        return

    # Criar credenciais SSL (Criptografia e Confiança do Servidor)
    credentials = grpc.ssl_channel_credentials(root_certs)
    
    # 2. Configurar a Conexão Segura
    with grpc.secure_channel(SERVER_ADDRESS, credentials) as channel:
        
        # 3. Criar o Stub (Cliente)
        stub = simple_pb2_grpc.OrderServiceStub(channel)

        # 4. Configurar Metadata (Token de Autenticação)
        metadata = [('authorization', AUTH_TOKEN)]

        # 5. Criar o Payload
        request = simple_pb2.OrderRequest(
            product_id='PYTHON-SERVICE-ITEM-42',
            quantity=1,
            price=99.99
        )
        
        print(f"Enviando solicitação OrderRequest para {SERVER_ADDRESS}...")
        
        try:
            # 6. Chamar o RPC
            response = stub.ProcessOrder(request, metadata=metadata)

            # 7. Processar a Resposta
            print("\n--- RESPOSTA BEM-SUCEDIDA (Python Client) ---")
            print(f"ID do Pedido: {response.order_id}")
            print(f"Total: R$ {response.total_amount:.2f}")
            print(f"Status: {response.status}")
            print("---------------------------------------------")

        except grpc.RpcError as e:
            # 8. Tratar Erros gRPC
            print("\n--- ERRO RPC ---")
            print(f"Status Code: {e.code()}")
            print(f"Detalhe: {e.details()}")
            print("----------------")

if __name__ == '__main__':
    run()