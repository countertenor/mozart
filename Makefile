build-all: ui add-static clean build build-linux build-centos
build:
	go build -ldflags "-X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-darwin main.go
build-w-clean: clean build
build-linux: # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/countertenor/mozart/internal/command.stateDBPathFromEnv=/tmp -X github.com/countertenor/mozart/internal/command.logDirPathFromEnv=/var/log/mozart -X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-linux main.go
build-centos: # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=ppc64le go build -ldflags "-X github.com/countertenor/mozart/internal/command.stateDBPathFromEnv=/tmp -X github.com/countertenor/mozart/internal/command.logDirPathFromEnv=/var/log/mozart -X github.com/countertenor/mozart/cmd.gitCommitHash=`git rev-parse HEAD` -X github.com/countertenor/mozart/cmd.buildTime=`date -u '+%Y-%m-%d--%H:%M:%S%p'` -X github.com/countertenor/mozart/cmd.gitBranch=`git branch --show-current`" -o bin/mozart-centos main.go
clean:
	rm -f bin/*
	rm -rf generated
	rm -rf logs
	rm -f *.db
	rm -f *.log
npm-install:
	(cd webapp; npm install)
ui: npm-install
	(cd webapp; npm run build)
install: add-static clean
	go install
add-static: #add static code to binary. if error: do 'go get github.com/rakyll/statik'
	statik -src resources -ns template -dest=statik/tmpl -f
	statik -src webapp/build -ns webapp -dest=statik/web -f
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