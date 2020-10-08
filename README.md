# Mozart

An orchestrator for scripts.

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Commands](#commands)
- [Building the application](#building-the-application)
- [Helpful resources](#helpful-resources)
  - [Good links for templating](#good-links-for-templating)
  - [Bash scripts through go](#bash-scripts-through-go)
  - [Versioning with go](#versioning-with-go)
  - [Web sockets](#web-sockets)
  - [Command execution](#command-execution)

<!-- /code_chunk_output -->

## Commands

```
mozart commands

- init          (creates yaml config file)
- parse         (optional, parses config file for error checking)
- execute       (executes all scripts in specified directory)
- state         (displays install state of all components, accepts optional args)
- version       (displays version info for the application)
  --json                (gives output in JSON)

Global flags
  -c           (configuration file, defaults to 'mozart-sample.yaml')
  -v           (prints verbosely, useful for debugging)
  -d, --dry-run       (optional) shows what scripts will run, but does not run the scripts
  -n, --no-generate   (optional) do not generate bash scripts as part of install/cleanup, instead use the ones in generated folder. Useful for running local change to the scripts
  -r, --re-run        (optional) re-run script from initial state, ignoring previously saved state

```

## Building the application

1. Clone the repo
1. `make install`
1. Make sure your `PATH` env variable has `go/bin` included (`export PATH=$PATH:$(go env GOPATH)/bin)`)

## Helpful resources

### Good links for templating

- https://golang.org/pkg/text/template/
- https://forum.golangbridge.org/t/referencing-map-values-in-text-template/6253/5
- https://medium.com/@IndianGuru/understanding-go-s-template-package-c5307758fab0
- https://goinbigdata.com/example-of-using-templates-in-golang/
- https://stackoverflow.com/questions/25689829/arithmetic-in-go-templates
- https://stackoverflow.com/questions/21305865/golang-separating-items-with-comma-in-template
- https://helm.sh/docs/chart_template_guide/functions_and_pipelines/

### Bash scripts through go

- https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
- https://stackoverflow.com/questions/18986943/in-golang-how-can-i-write-the-stdout-of-an-exec-cmd-to-a-file
- https://stackoverflow.com/questions/19965795/how-to-write-log-to-file
- https://yourbasic.org/golang/format-parse-string-time-date-example/

### Versioning with go

- https://www.atatus.com/blog/golang-auto-build-versioning/

### Web sockets

- https://www.zupzup.org/io-pipe-go/
- https://gist.github.com/ifels/10392762
- https://godoc.org/github.com/gorilla/websocket

### Command execution

- https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
- https://stackoverflow.com/a/58572436

### Error stack - todo

- https://github.com/gruntwork-io/terratest/blob/master/modules/helm/template.go#L83
- https://godoc.org/github.com/go-errors/errors
- https://github.com/gruntwork-io/gruntwork-cli/blob/master/errors/errors.go

### Embedding

- https://github.com/golang/go/issues/41191
