echo "Genetating proto files..."

protoc --go_out=./stream --go_opt=paths=source_relative \
    --go-grpc_out=./stream --go-grpc_opt=paths=source_relative \
    ./stream.proto
