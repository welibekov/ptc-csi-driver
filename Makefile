# Makefile

# Define the binary name and source file
BINARY = bin/ptc-csi-driver
SRC = cmd/main.go

# Define the Go command
GO = go

# Default rule
all: $(BINARY)

# Build rule
$(BINARY): $(SRC)
	@mkdir -p bin
	$(GO) build -o $@ $<

# Clean rule
clean:
	rm -f $(BINARY)
