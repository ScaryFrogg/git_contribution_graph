build-local:
	@echo "Building... to /usr/local/bin"
	@go build -o /usr/local/bin/gcg cmd/git_contribution_graph/main.go
run:
	@go run cmd/git_contribution_graph/main.go
	@echo "$(CURRENT_VERSION_MICRO)"
clean:
	@echo "Cleaning..."
	@rm ./bin/*
test:
	@go run cmd/git_contribution_graph/main.go -token=$$GH_CONTRIBUTION_KEY -username=ScaryFrogg
	@echo ""
	@go run cmd/git_contribution_graph/main.go -token=$$GH_CONTRIBUTION_KEY -username=ScaryFrogg -from=2024-05-23T00:00:00Z
	@echo ""
	@go run cmd/git_contribution_graph/main.go -token=$$GH_CONTRIBUTION_KEY -username=ScaryFrogg -from=2024-05-23T00:00:00Z -colors=1,230,11,157,10
	@echo ""
	@go run cmd/git_contribution_graph/main.go -token=$$GH_CONTRIBUTION_KEY -username=ScaryFrogg -from=2024-05-23T00:00:00Z -legend=false

DESCRIBE           := $(shell git describe --match "*" --always --tags)
DESCRIBE_PARTS     := $(subst -, ,$(DESCRIBE))

VERSION_TAG        := $(word 1,$(DESCRIBE_PARTS))
COMMITS_SINCE_TAG  := $(word 2,$(DESCRIBE_PARTS))

VERSION            := $(subst v,,$(VERSION_TAG))
VERSION_PARTS      := $(subst ., ,$(VERSION))

MAJOR              := $(word 1,$(VERSION_PARTS))
MINOR              := $(word 2,$(VERSION_PARTS))
MICRO              := $(word 3,$(VERSION_PARTS))

NEXT_MAJOR         := $(shell echo $$(($(MAJOR)+1)))
NEXT_MINOR         := $(shell echo $$(($(MINOR)+1)))
NEXT_MICRO	    = $(shell echo $$(($(MICRO)+$(COMMITS_SINCE_TAG))))

ifeq ($(strip $(COMMITS_SINCE_TAG)),)
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(MICRO)
CURRENT_VERSION_MINOR := $(CURRENT_VERSION_MICRO)
CURRENT_VERSION_MAJOR := $(CURRENT_VERSION_MICRO)
else
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(NEXT_MICRO)
CURRENT_VERSION_MINOR := $(MAJOR).$(NEXT_MINOR).0
CURRENT_VERSION_MAJOR := $(NEXT_MAJOR).0.0
endif

build-release:
	@echo "Building version $(CURRENT_VERSION_MICRO)..."
	@echo "Building Linux binary..."
	@go build -o ./bin/gcg cmd/git_contribution_graph/main.go
	@echo "Packaging Linux binary..."
	@tar -czvf ./bin/linux_amd64.tar.gz -C ./bin/ gcg

	@echo "Building Windows binary..."
	@env GOOS=windows GOARCH=amd64 go build -o ./bin/gcg.exe cmd/git_contribution_graph/main.go
	@echo "Packaging Windows binary..."
	@zip -j ./bin/windows_amd64.zip ./bin/gcg.exe

release-micro:
	@echo "Building version $(CURRENT_VERSION_MICRO)"
	@$(MAKE) build-release

	@echo "Publishing version $(CURRENT_VERSION_MICRO)"
	@git tag -a $(CURRENT_VERSION_MICRO) -m "Release version $(CURRENT_VERSION_MICRO)"
	@git push origin $(CURRENT_VERSION_MICRO)

release-minor:
	@echo "Building version $(CURRENT_VERSION_MINOR)"
	@$(MAKE) build-release

	@echo "Publishing version $(CURRENT_VERSION_MINOR)"
	@git tag -a $(CURRENT_VERSION_MINOR) -m "Release version $(CURRENT_VERSION_MINOR)"
	@git push origin $(CURRENT_VERSION_MINOR)

release-major:
	@echo "Building version $(CURRENT_VERSION_MAJOR)"
	@$(MAKE) build-release

	@echo "Publishing version $(CURRENT_VERSION_MAJOR)"
	@git tag -a $(CURRENT_VERSION_MAJOR) -m "Release version $(CURRENT_VERSION_MAJOR)"
	@git push origin $(CURRENT_VERSION_MAJOR)

.PHONY: build-local run clean build-release release-major release-micro release-minor test
