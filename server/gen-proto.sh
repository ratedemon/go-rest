
DIRECTORIES=$(cd proto && ls -d -- */)

for proto_dir in $DIRECTORIES
do
  proto_file=$(find ./proto/$proto_dir -type f -name "*.proto" | head -n 1)
  protoc -I . \
    -I ${GOPATH}/src \
    -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
    ./$proto_file --go_out=plugins=grpc:. \
    --validate_out="lang=go:."
  echo "Proto is generated - ./$proto_dir"
done