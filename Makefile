.PHONY: run build clean

BINARY := scraper
GOENV  := GOTOOLCHAIN=local CGO_ENABLED=0

# Build and run the scraper interactively
run: build
	@./$(BINARY)

# Compile the binary
build:
	@echo "Building..."
	@$(GOENV) go build -o $(BINARY) .
	@echo "Done. Run with: ./$(BINARY)"

# Remove the compiled binary
clean:
	@rm -f $(BINARY)
	@echo "Cleaned."
