PROTO_DIR = proto
GEN_DIR = gen/go

.PHONY: proto
proto:
	@echo "генерация protobuf..."
	@mkdir -p $(GEN_DIR)
	
	protoc -I $(PROTO_DIR) \
	-I third_party/googleapis \
	--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=paths=source_relative \
	$(PROTO_DIR)/auth/auth_service.proto \
	$(PROTO_DIR)/user/user_service.proto \
	$(PROTO_DIR)/task/task_service.proto \
	$(PROTO_DIR)/group/group_service.proto

.PHONY: clean
clean:
	rm -rf $(GEN_DIR)