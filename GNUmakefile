TEST?=$$(go list ./... | grep -v 'vendor' | grep -v 'env' | grep -v 'utils')
NAME=ya41-56
VERSION=0.0.1
BINARY_GOPHERMART=gophermart
BINARY_ACCRUAL=accrual
OS_ARCH=darwin_amd64
COMBINED_FILE := combined.go

default: build

build: build_gophermart build_accrual

combine_mart:
	@echo "Combining Go files into $(COMBINED_FILE)..."
	@rm -f $(COMBINED_FILE)
	@{ \
		echo cmd/flags.go; \
		find cmd/gophermart -type f -name "*.go"; \
		find internal/gophermart -type f -name "*.go"; \
		find internal/shared -type f -name "*.go"; \
	} | while read -r file; do \
		echo "// ===== $$file =====" >> $(COMBINED_FILE); \
		cat "$$file" >> $(COMBINED_FILE); \
		echo "" >> $(COMBINED_FILE); \
	done
	@echo "Done: $(COMBINED_FILE)"

combine_accrual:
	@echo "Combining Go files into $(COMBINED_FILE)..."
	@rm -f $(COMBINED_FILE)
	@{ \
		echo cmd/flags.go; \
		find cmd/accrual -type f -name "*.go"; \
		find internal/accrual -type f -name "*.go"; \
		find internal/shared -type f -name "*.go"; \
	} | while read -r file; do \
		echo "// ===== $$file =====" >> $(COMBINED_FILE); \
		cat "$$file" >> $(COMBINED_FILE); \
		echo "" >> $(COMBINED_FILE); \
	done
	@echo "Done: $(COMBINED_FILE)"

build_accrual:
	go build -o ${BINARY_ACCRUAL} cmd/accrual/main.go

build_gophermart:
	go build -o ${BINARY_GOPHERMART} cmd/gophermart/main.go

test_gophermart:
	go test ./... -timeout 60m --tags=gophermart -v

test_accrual:
	go test ./... -timeout 60m --tags=accrual -v

run_test_mart:
	./gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=./gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgres://myuser:mypassword@localhost:5432/mydatabase" \
		-accrual-binary-path=./accrual \
		-accrual-host=localhost \
		-accrual-port=41222 \
		-accrual-database-uri="postgres://myuser:mypassword@localhost:5432/mydatabase"

run_test_mart_original:
	./gophermarttest-darwin-arm64 \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=./gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable" \
		-accrual-binary-path=./accrual \
		-accrual-host=localhost \
		-accrual-port=41222 \
		-accrual-database-uri="postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

lint:
	golangci-lint run ./...