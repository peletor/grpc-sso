#! /bin/bash
protoc -I proto proto/sso/sso.proto --go_out=./gen/ --go_opt=paths=source_relative --go-grpc_out=./gen/ --go-grpc_opt=paths=source_relative
if [ $? -eq 0 ]; then
  echo "Files were generated successfully"
fi