build:
	go build -ldflags "-X github.com/prashantgupta24/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/prashantgupta24/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/prashantgupta24/mozart/cmd.gitBranch=`git branch --show-current` -X github.com/prashantgupta24/mozart/cmd.tagVersion=`git describe --tags --long`" -o bin/mozart main.go
build-w-clean: clean build
build-linux: # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/prashantgupta24/mozart/internal/command.stateDBPathFromEnv=/tmp -X github.com/prashantgupta24/mozart/internal/command.logDirPathFromEnv=/var/log/mozart -X github.com/prashantgupta24/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/prashantgupta24/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/prashantgupta24/mozart/cmd.gitBranch=`git branch --show-current` -X github.com/prashantgupta24/mozart/cmd.tagVersion=`git describe --tags --long`" -o bin/mozart main.go
clean:
	rm -f bin/mozart
	rm -rf generated
	rm -rf logs
	rm -f *.db
	rm -f *.log
install: add-static clean
	go install
add-static: #add static code to binary. if error: do 'go get github.com/rakyll/statik'
	statik -src static -f
run-help:
	go run main.go --help
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
start-ws:
	go run ws/ws.go