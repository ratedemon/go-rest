
DIRECTORIES=$(cd proto && ls -d -- */)

for proto_dir in $DIRECTORIES
do
  proto_file=$(find ./proto/$proto_dir -type f -name "*.proto" | head -n 1)
  protoc -I . ./$proto_file --go_out=plugins=grpc:.
  echo "Proto is generated - ./$proto_dir"
done