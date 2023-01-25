BINARY_NAME=study-redis-cluster
 
build:
	go build -o ${BINARY_NAME} -ldflags="-s -w" ./cmd/${BINARY_NAME}
 
clean:
	go clean
	rm ${BINARY_NAME}
 
tc:
	go test ./...
 
test_coverage:
	go test ./... -coverprofile=coverage.out
 
dep:
	go mod download
 
run:
	./${BINARY_NAME}
 
all: build
