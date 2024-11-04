# Variables

GO = go
PROJECT_DIR = ${shell pwd}

# Environment variables

# Commands


build:
	@echo "==> (1/2) Building the project..."
	@${GO} mod tidy
	@${GO} build -o ${PROJECT_DIR}/cmd/bin/server ${PROJECT_DIR}/cmd/server/main.go
	@echo "==> (2/2) Project built successfully!"

run: build
	@echo "==> (1/1) Running the project..."
	@${PROJECT_DIR}/cmd/bin/server

generate:
	@echo "==> (1/2) Generating docs code..."
	@swag init -g cmd/server/main.go 
	@echo "==> (2/2) Docs code generated successfully!"

clean:
	@echo "==> (1/2) Cleaning the project..."
	@rm -rf ${PROJECT_DIR}/cmd/bin
	@echo "==> (2/2) Project cleaned successfully!"

.DEFAULT_GOAL := run