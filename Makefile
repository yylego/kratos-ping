PING_PROTO_FILES = $(shell find ./clientpingkratos -name "*.proto")

.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@echo "All protoc plugins installed!"

.PHONY: generate
generate:
	protoc --proto_path=./clientpingkratos \
		   --proto_path=./internal/proto3ps \
		   --go_out=paths=source_relative:./clientpingkratos \
		   --go-http_out=paths=source_relative:./clientpingkratos \
		   --go-grpc_out=paths=source_relative:./clientpingkratos \
		   --openapi_out=fq_schema_naming=true,default_response=false,title=PING-KRATOS:. \
		   $(PING_PROTO_FILES)
	@echo "Proto code generation complete!"

# ========================================
# TEMPLATE BEGIN: TEST AND COVERAGE CONFIG
# ========================================
# Test and Coverage (GitHub Actions)
# ========================================

COVERAGE_DIR ?= .coverage.out

test:
	@if [ -d $(COVERAGE_DIR) ]; then rm -r $(COVERAGE_DIR); fi
	@mkdir $(COVERAGE_DIR)
	make test-with-flags TEST_FLAGS='-v -race -covermode atomic -coverprofile $$(COVERAGE_DIR)/combined.txt -bench=. -benchmem -timeout 20m'

test-with-flags:
	@go test $(TEST_FLAGS) ./...

# ========================================
# TEMPLATE END: TEST AND COVERAGE CONFIG
# ========================================

.PHONY: clean
clean:
	rm -f ./clientpingkratos/*.pb.go
	rm -f ./openapi.yaml
	@echo "Cleanup complete!"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  install  - Install protoc plugins (go, grpc, http, openapi)"
	@echo "  generate - Generate Go code from proto files"
	@echo "  test     - Run tests with coverage"
	@echo "  clean    - Remove generated files"
	@echo "  help     - Show this help message"
