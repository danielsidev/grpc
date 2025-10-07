const PROTO_PATH = './simple.proto';
const CERT_PATH = './server.pem'; // <-- Caminho para o certificado público
const fs = require('fs');         // <-- Módulo para ler o arquivo
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

// --- Configuração ---
const SERVER_ADDRESS = 'localhost:50051';
const AUTH_TOKEN = 'Bearer my-secret-api-key-12345'; 

// 1. Carregar o Protobuf (inalterado)
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
const simpleProto = grpc.loadPackageDefinition(packageDefinition).simple;


// 2. Definir Credenciais (TLS Seguro)
let credentials;
try {
    // Carregar o certificado público do servidor (CA Raiz)
    const rootCert = fs.readFileSync(CERT_PATH); 
    
    // Criar credenciais SSL usando o certificado.
    // O primeiro argumento (rootCert) é a CA Raiz que o cliente confia.
    credentials = grpc.credentials.createSsl(rootCert);
    
    console.log("Credenciais TLS carregadas com sucesso a partir de server.pem.");
} catch (e) {
    console.error(`Erro ao carregar o certificado em ${CERT_PATH}. Certifique-se de que o arquivo existe.`);
    console.error(e.message);
    process.exit(1);
}


// 3. Criar o Cliente (Stub)
// A conexão agora é forçada a ser segura (TLS) e verifica a identidade do servidor
const client = new simpleProto.OrderService(
    SERVER_ADDRESS, 
    credentials 
);

// 4. Implementar a chamada RPC (inalterado)
function runProcessOrder() {
    // Payload de ENTRADA (OrderRequest)
    const request = {
        product_id: 'NODE-SECURE-ITEM-1',
        quantity: 10,
        price: 7.99
    };

    // Metadata para Autenticação (Token)
    const metadata = new grpc.Metadata();
    metadata.add('authorization', AUTH_TOKEN);

    console.log(`\nEnviando solicitação segura para ${SERVER_ADDRESS}...`);
    
    // Fazer a chamada RPC
    client.processOrder(request, metadata, (error, response) => {
        if (error) {
            // Em caso de erro TLS, o código será 14 (UNAVAILABLE) ou 13 (INTERNAL)
            // Em caso de erro de Token, o código será 16 (UNAUTHENTICATED)
            console.error('\n--- ERRO RPC ---');
            console.error(`Status Code: ${error.code}`); 
            console.error(`Detalhe: ${error.details}`);
            console.error('----------------');
            return;
        }

        // Payload de SAÍDA (ReceiptResponse)
        console.log('\n--- RESPOSTA BEM-SUCEDIDA (CONEXÃO SEGURA) ---');
        console.log(`Recebido ID do Pedido: ${response.order_id}`);
        console.log(`Total: R$ ${response.total_amount.toFixed(2)}`);
        console.log(`Status: ${response.status}`);
        console.log('---------------------------------------------');
    });
}

// Executar a função
runProcessOrder();