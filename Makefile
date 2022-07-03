build:
	go build -o gorevershserver server.go
	go build -o gorevershclient client.go

runserver:
	go run server.go

runclient:
	go run client.go