.DEFAULT_GOAL := help
SHELL := /bin/bash

#help: @ list available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#test.unit: @ run unit tests and coverage
test.unit:
	@echo "[TEST.UNIT] run unit tests and coverage"
	@go test -race -covermode=atomic -coverprofile=coverage.out \
		github.com/SLedunois/b3lbctl/pkg/cmd/root \
		github.com/SLedunois/b3lbctl/pkg/cmd/instances \
		github.com/SLedunois/b3lbctl/pkg/cmd/clusterinfo \
		github.com/SLedunois/b3lbctl/pkg/cmd/config \
		github.com/SLedunois/b3lbctl/pkg/admin \
		github.com/SLedunois/b3lbctl/pkg/config \
		github.com/SLedunois/b3lbctl/pkg/system