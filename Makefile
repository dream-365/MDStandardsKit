BINARY := mdstandards
GO ?= go

.PHONY: build test vet check fmt clean

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

check: build vet test
	./$(BINARY) -root . -policy ./policy -phase 1 -format table

fmt:
	$(GO) fmt ./...

clean:
	rm -f $(BINARY)
