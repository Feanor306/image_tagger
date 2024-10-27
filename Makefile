# Build commands
templ-install:
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if [ ! -x "$$(command -v templ)" ]; then \
				echo "templ installation failed. Exiting..."; \
				exit 1; \
			fi; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi

build: templ-install
	@echo "Building..."
	@templ generate
	
	@go build -o bin/main cmd/api/main.go

# Run the application
run: build
	@./bin/main

# installs godoc and runs local doc server at http://localhost:6060
doc:
	go install golang.org/x/tools/cmd/godoc@latest 
	godoc -http=:6060 -goroot=.

# Create DB container
db-up:
	@if docker compose -f dockerized/database.yml up --build -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f dockerized/database.yml up --build -d; \
	fi

# Shutdown DB container
db-down:
	@if docker compose -f dockerized/database.yml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f dockerized/database.yml down; \
	fi

# Create test DB container
test-db-up:
	@if docker compose -f dockerized/db_test.yml up --build -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f dockerized/db_test.yml up --build -d; \
	fi

test-db-wait:
	@echo "Creating test db..." 
	sleep 10

# Shutdown test DB container
test-db-down:
	@if docker compose -f dockerized/db_test.yml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f dockerized/db_test.yml down; \
	fi

run-tests: 
	@if go test ./src/... 2>/dev/null; then \
		: ; \
	fi

test: test-db-up test-db-wait run-tests test-db-down
