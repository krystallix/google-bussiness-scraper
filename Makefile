.PHONY: run build serve filter clean

BINARY := scraper
GOENV  := GOTOOLCHAIN=local CGO_ENABLED=0

# Build and run the scraper interactively (CLI mode)
run: build
	@./$(BINARY)

# Compile the server binary
build:
	@echo "Building server..."
	@$(GOENV) go build -o $(BINARY) ./cmd/server/
	@echo "Done. Run with: ./$(BINARY)"

# Start the web GUI dashboard
serve:
	@echo "Starting web dashboard at http://127.0.0.1:8080"
	@$(GOENV) go run ./cmd/server/

# Filter and classify existing CSV output into leads report
filter:
	@echo "Running lead classifier..."
	@$(GOENV) go run ./cmd/filter/

# Remove the compiled binary
clean:
	@rm -f $(BINARY)
	@echo "Cleaned."
