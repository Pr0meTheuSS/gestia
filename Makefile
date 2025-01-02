
.PHONY: run
run:
	go run cmd/main.go

.PHONY: generate-coverage
generate-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: generate-swagger
generate-swagger:
	swag init -d "./" -g "cmd/main.go"