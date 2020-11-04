# Mozart

An orchestrator for scripts.

Mozart is a simple utility to convert your independent, messy bunch of scripts into a well defined, orchestrated program!

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Easy steps to migrate to Mozart](#easy-steps-to-migrate-to-mozart)
- [What exactly does an orchestrator do?](#what-exactly-does-an-orchestrator-do)
- [Benefits of mozart](#benefits-of-mozart)
- [Details](#details)
- [Commands](#commands)
- [Building the application](#building-the-application)
- [Helpful resources](#helpful-resources)
  - [Good links for templating](#good-links-for-templating)
  - [Bash scripts through go](#bash-scripts-through-go)
  - [Versioning with go](#versioning-with-go)
  - [Web sockets](#web-sockets)
  - [Command execution](#command-execution)
  - [Error stack - todo](#error-stack-todo)
  - [Embedding static content](#embedding-static-content)

<!-- /code_chunk_output -->

## Easy steps to migrate to Mozart

It is extremely simple - NO CODE CHANGE REQUIRED!

1. Dump your bash or python scripts into the `resources/templates` directory. (Modularize if necessary. Use templating if necessary.)
1. Build the binary.
1. Voila! You have a single orchestrator binary with all your scripts in it!

## What exactly does an orchestrator do?

Good question. Simply speaking, an orchestrator manages execution of your unorganized scripts. Suppose you had 2 bash scripts which let you install and uninstall a particular component respectively.

Without an orchestrator, :

1. You would have to ship both files to anyone who wants to use them.
1. If running on a shared system, you can never know if the other person already ran the install script (unless actually going through the effort of seeing if the component was actually installed).
1. No way of using common variables across the scripts (like version of the component being managed).
1. No way of accidentally preventing execution of a script more than once (like prevent install script from running again if it already ran once).
1. Manually execute the scripts through bash or python (no CLI or UI)

## Benefits of mozart

- Simple migration - no code change necessary.
- Lets you modularize the scripts, which means you can have smaller scripts which do smaller tasks (called modules). No need to maintain huge bash files anymore. The smaller the scripts - the easier it is to manage them.
- Ability to use templating capabilities - similar to helm. Values in the yaml file are read by all scripts.
- Public state of execution of scripts, so no more confusion of whether scripts were executed or not.
- Make sure scripts are executed just once - the orchestrator will not allow you to run the same script again without explicitely mentioning a `Re-Run` flag, thereby preventing accidental execution of the same script.
- Ready to use CLI + UI to manage the execution.
- Single binary which contains everything - no need to send in bunch of scripts to anyone anymore.

## Details

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

### Embedding static content

- https://github.com/rakyll/statik/pull/101/files
- https://github.com/golang/go/issues/41191
