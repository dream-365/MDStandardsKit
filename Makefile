BINARY := mdstandards
TASKSCHECK := taskscheck
GO ?= go

.PHONY: build test vet check taskscheck fmt clean

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

check: build vet test taskscheck
	./$(BINARY) -root . -policy ./policy -phase 1 -format table
	./$(TASKSCHECK) -root . -phase 1

taskscheck:
	$(GO) build -o $(TASKSCHECK) ./tools/taskscheck

fmt:
	$(GO) fmt ./...

clean:
	rm -f $(BINARY) $(TASKSCHECK)
