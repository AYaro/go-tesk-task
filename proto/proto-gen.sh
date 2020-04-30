protoc --proto_path=api/proto/v1 --proto_path=proto --go_out=plugins=grpc:pkg/api/v1 parts-service.proto
