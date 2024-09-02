PROTO_NAMES=(
    "model"
    "gateway"
)

for name in "${PROTO_NAMES[@]}"; do
  protoc --go_out=paths=source_relative:protocol --go-grpc_out=paths=source_relative:protocol --proto_path=protocol protocol/${name}/${name}.proto
  if [ $? -ne 0 ]; then
      echo "error processing ${name}.proto"
      exit $?
  fi
done

if [ "$(uname -s)" == "Darwin" ]; then
    find . -type f -name '*.pb.go' -exec sed -i '' 's/,omitempty"`/\"\`/g' {} +
else
    find . -type f -name '*.pb.go' -exec sed -i 's/,omitempty"`/\"\`/g' {} +
fi
