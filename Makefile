
codegen:  generate

revendor:
	@go mod tidy -v
	@go mod vendor -v
	@go mod verify

generate:
	echo "Building $@"
	./generate_grpc.sh

.PHONY:
