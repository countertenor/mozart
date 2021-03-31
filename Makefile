BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

GIT_ADD=github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`
DB_PATH=github.com/countertenor/mozart/internal/command.stateDBPathFromEnv
LOG_PATH=github.com/countertenor/mozart/internal/command.logDirPathFromEnv

build-all: ui build-all-binary
build-all-binary: build-darwin build-linux build-centos
build-darwin: clean
	go build -ldflags "-X $(GIT_ADD)" -o bin/darwin/mozart main.go
build-linux: clean
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X $(DB_PATH)=/tmp -X $(LOG_PATH)=/var/log/mozart -X $(GIT_ADD)" -o bin/linux/mozart main.go
build-centos: clean
	env GOOS=linux GOARCH=ppc64le go build -ldflags "-X $(DB_PATH)=/tmp -X $(LOG_PATH)=/var/log/mozart -X $(GIT_ADD)" -o bin/centos/mozart main.go
clean:
	rm -f bin/*.*
	rm -rf generated
	rm -rf logs
	rm -f *.db
	rm -f *.log
npm-install:
	(cd static/webapp; npm install)
ui: npm-install
	(cd static/webapp; npm run build)
install: clean
	go install -tags ui
install-no-ui: clean
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