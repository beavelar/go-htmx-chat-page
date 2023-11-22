# Go/HTMX Chat Page

## Compiling protobuf for Go
protoc --proto_path=proto --go_out=services/content/genproto database.proto --go-grpc_out=services/content/genproto database.proto
protoc --proto_path=proto --go_out=services/database/genproto database.proto --go-grpc_out=services/database/genproto database.proto

## Compiling templ files
templ generate
