.PHONY: run build filter clean

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

# Filter and classify existing CSV output into leads report
filter:
	@echo "Running lead classifier..."
	@$(GOENV) go run ./cmd/filter/

# Remove the compiled binary
clean:
	@rm -f $(BINARY)
	@echo "Cleaned."
