# Guidelines 

> In progress....wait for more...

---

```
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       simple.proto
```


```

# 1. Gerar a chave privada do servidor (server.key)
openssl genrsa -out server.key 2048

# 2. Gerar o certificado (server.pem) usando o arquivo de configuração
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 365 -config cert.conf -extensions v3_req
```

```
python3 -m venv venv_grpc
source venv_grpc/bin/activate
pip install -r requirements.txt

python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. simple.proto
```