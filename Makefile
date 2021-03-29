build-all: ui build-all-binary
build-all-binary: build-darwin build-linux build-centos
build-darwin: clean
	go build -ldflags "-X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-darwin main.go
build-linux: clean # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/countertenor/mozart/internal/command.stateDBPathFromEnv=/tmp -X github.com/countertenor/mozart/internal/command.logDirPathFromEnv=/var/log/mozart -X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-linux main.go
build-centos: clean # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=ppc64le go build -ldflags "-X github.com/countertenor/mozart/internal/command.stateDBPathFromEnv=/tmp -X github.com/countertenor/mozart/internal/command.logDirPathFromEnv=/var/log/mozart -X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-centos main.go
clean:
	rm -f bin/*
	rm -rf mozart-generated
	rm -rf logs
	rm -f *.db
	rm -f *.log
npm-install:
	(cd static/webapp; npm install)
ui: npm-install
	(cd static/webapp; npm run build)
install: clean
	go install
run-server: install
	mozart server
server-live: # go get -u github.com/cosmtrek/air
	air -c .air.toml
# save github token in an environment variable export GITHUB_TOKEN="token"
add-tag:
ifeq (,$(findstring v,$(tag)))
	@echo "error : tag needs to be of format v0.x.x. Usage --> make upload tag=v0.x.x"
	@echo
	exit 1
endif
	git fetch
	git tag $(tag)
	git push origin --tags
upload: add-tag install build-linux #make upload tag=v0.x.x, install --> brew install goreleaser
	goreleaser --rm-dist
test:
	go test -v ./...