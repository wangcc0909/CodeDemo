build:
    protoc -I. --go_out=plugins=micro:($GOPATH)/src/shippy2/consignment-server ($GOPATH)/src/shippy2/consignment-server/consignment.proto
    GOOS=linux GOARCH=amd64 go build
    docker build -t consignment-server .
run:
    docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-server