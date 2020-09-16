default: help


BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)


# set target color
TARGET_COLOR := $(BLUE)

.PHONY: help
help: ## - Show help message
	@printf "${TARGET_COLOR} usage: make [target]\n${RESET}"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s${RESET} %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

lint: ## - Linter
	@echo "${TARGET_COLOR} Lint code !${RESET}" ;\
	go vet ./...
