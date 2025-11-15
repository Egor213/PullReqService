#!/bin/sh

mockgen -source=internal/repo/repo.go \
        -destination=internal/mocks/repomock/repomocks.go \
        -package=repomocks

mockgen -source=internal/service/service.go \
        -destination=internal/mocks/servicemock/servicemocks.go \
        -package=servmocks
