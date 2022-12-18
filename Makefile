BINDIR:=bin

ROOT_PACKAGE=contentssecurity
COMMAND_PACKAGES:=$(shell go list ./cmd/...)

BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%.o)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

.PHONY: build
build: $(BINARIES)

$(BINARIES): $(GO_FILES)
	@mkdir -p $(BINDIR)
	go build -o $@ $(@:$(BINDIR)/%.o=$(ROOT_PACKAGE)/cmd/%)

.PHONY: clean
clean:
	@rm -rf $(BINDIR)/*