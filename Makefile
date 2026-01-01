
dev:
	@if command -v air > /dev/null; then \
            air; \
    elif [ -f "$$(go env GOPATH)/bin/air" ]; then \
            "$$(go env GOPATH)/bin/air"; \
    else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                "$$(go env GOPATH)/bin/air"; \
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
    fi

run:
	go run cmd/server/main.go
