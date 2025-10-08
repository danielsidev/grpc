# Guidelines 

> In progress....wait for more...
---

#### Intro

To create a GRPC communication, we need , first, of a proto file that contain the rules of our application. Based in proto file, in this case,  simple.proto, we generate new files to start our application.

### Proto Buffers Files 

In this example, we already have the files, but in a new application, we need to generate this files.

Generating proto buffers files for a Server Application Golang GRPC:

```
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       simple.proto
```

### Authentication

How our application's example  works with security in two layers, we need generate a private and public key. Our server application will user a private key and our client applicationns will use a public key.

1. First, generate the private key to server application ( **server.key** ) :

```
openssl genrsa -out server.key 2048

```
> The private key must be protect and only your server application can access


2. Second, generate the public key to client applications ( **server.pem** ) :

> **This is an example to a self-signed certificate. Because this we use the cert.conf file.**

```
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 365 -config cert.conf -extensions v3_req

```

> Share the public key with your client applications that can user your services ( for this example, copy the server.pem to client applications folders )


### Authorization

In addition to authentication, it is good practice to have authorization in our application.

For this, we use an AUTH_TOKEN, similar to API_KEY to control access.

Generate a enviroment variable called AUTH_TOKEN  and define your token.

Example:

```
export AUTH_TOKEN='Bearer my-secret-api-key-12345'
```

> So, in addition from server.pem, our client applications must be to send this token in header to have a authorizatition confirmed.

Examples:

In Golang Client

```
	"... 
       
       grpc.WithPerRPCCredentials(auth) 
       
    ..."
```
In NodeJs Client

```
   "... 
    
    const metadata = new grpc.Metadata();
    metadata.add('authorization', AUTH_TOKEN); 
    
    ..."
```

In Python Client

```
   "...
   
    stub = simple_pb2_grpc.OrderServiceStub(channel)
    metadata = [('authorization', AUTH_TOKEN)]
    ...
   
    response = stub.ProcessOrder(request, metadata=metadata)

   ..."


```



### About simple.proto file in client applications


In this example, simple.proto is our configuration file for GRPC, so, it is that define our rules for our project, like payloads, responses object and  availables services.

Because this, all our client application need a copy this file. 
Our client application use this file to  know which services and contracts are availbale in server application.

### How to Run this Project

We have a server application in Golang, so, first, we need:

#### Golang Server Application

- Install depencies from Golang App ( inside go_service folder in temrinal )

```
go mod tidy

```

- Gets Up the Server Application ( **inside go_service/server/**  )

```
go run main.go

```
#### Golang Client Application


Now in other terminal instance ( **inside go_service/client/** )

- Gets Up Golang's Client


```
go run main.go

```

#### NodeJs Client Application

In other terminal instance ( **inside node_service/** ):

- Install depencies from NodeJs App  

```
npm install

```
- Gets up the NodeJs' client
```
node clients.js

```
--- 

### Python Client Use Case

For python client application, we need generate some proto buffers files to communicate with GRPC server.

For this, to do ( **inside python_service folder in other terminal instance** ):

1. Remember to create your virtual environment

```
python3 -m venv venv_grpc
```


2. Active your virtual enviroment 

```
source venv_grpc/bin/activate

```

3.  Install the dependencies from python client

```
pip install -r requirements.txt

```

4. Generate the proto buffers files to python client

```

python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. simple.proto

```

5. Gets up the Python' Client

```
python client.py

```