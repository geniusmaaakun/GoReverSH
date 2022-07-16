build:
	go build -o goreversh_server server.go
	go build -o goreversh_client client.go

build_win:
	go build -o goreversh_server.exe server.go
	go build -o goreversh_client.exe client.go

runserver:
	go run server.go

runclient:
	go run client.go

test:
	go test ./config -v
	go test ./utils -v
	go test ./pkgclient -v
	go test ./pkgserver -v
