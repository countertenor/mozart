# Mozart

An orchestrator for scripts.

Mozart is a simple utility to convert your independent, messy bunch of scripts into a well defined, orchestrated program!

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Easy steps to migrate to Mozart](#easy-steps-to-migrate-to-mozart)
- [What exactly does an orchestrator do?](#what-exactly-does-an-orchestrator-do)
- [Benefits of using Mozart](#benefits-of-using-mozart)
- [How to add scripts to the orchestrator](#how-to-add-scripts-to-the-orchestrator)
  - [Steps](#steps)
    - [1. Initial setup](#1-initial-setup)
    - [2. (Optional) Modularize](#2-optional-modularize)
    - [3. (Optional) Use templating](#3-optional-use-templating)
    - [4. Build the binary](#4-build-the-binary)
- [CLI](#cli)
- [UI](#ui)
- [Helpful resources](#helpful-resources)
  - [Good links for templating](#good-links-for-templating)
  - [Bash scripts through go](#bash-scripts-through-go)
  - [Versioning with go](#versioning-with-go)
  - [Web sockets](#web-sockets)
  - [Command execution](#command-execution)
  - [Error stack](#error-stack)
  - [Embedding static content](#embedding-static-content)

<!-- /code_chunk_output -->

## Easy steps to migrate to Mozart

It is extremely simple - NO CODE CHANGE REQUIRED!

1. Create a directory inside `resources/templates`.
1. Dump your bash or python scripts into that directory. (Various optional steps can be included here - modularize if necessary, use templating if necessary. Discussed below)
1. Build the binary. (discussed below.)
1. Voila! You have a single orchestrator binary with all your scripts in it!

(discussed in detail below.)

## What exactly does an orchestrator do?

Good question. Simply speaking, an orchestrator manages execution of your unorganized scripts.

Suppose you had 2 bash scripts which let you install and uninstall a particular component respectively. Without an orchestrator:

1. You would have to ship both files to anyone who wants to use them.
1. If running on a shared system, you can never know if the other person already ran the install script (unless actually going through the effort of seeing if the component was actually installed).
1. No way of using common variables across the scripts (like version of the component being managed).
1. No way of accidentally preventing execution of a script more than once (like prevent install script from running again if it already ran once).
1. Manually execute the scripts through bash or python (no CLI or UI)

## Benefits of using Mozart

- Simple migration - no code change necessary.
- Lets you modularize the scripts, which means you can have smaller scripts which do smaller tasks (called modules). No need to maintain huge bash files anymore. The smaller the scripts - the easier it is to manage them.
- Ability to use templating capabilities - similar to helm. Values in the yaml file are read by all scripts.
- Public state of execution of scripts, so no more confusion of whether scripts were executed or not.
- Make sure scripts are executed just once - the orchestrator will not allow you to run the same script again without explicitely mentioning a `Re-Run` flag, thereby preventing accidental execution of the same script.
- Ready to use CLI + UI to manage the execution.
- Single binary which contains everything - no need to send in bunch of scripts to anyone anymore.

## How to add scripts to the orchestrator

Let us walk through how to actually add your scripts.

Once you clone the repo, the directory where you will add your scripts is going to be `resources/templates`.

**Modules**

The way to add your scripts is through `Modules`. Modules are nothing but directories, created with the intention of performing 1 simple task.

**Sample module**

There is already a sample module present there called `test-module`.

### Steps

#### 1. Initial setup

1. Clone the repo
1. Create a base directory (or module) which is the main directory under which all your scripts (or modules) will exist.

#### 2. (Optional) Modularize

1. Look through your scripts, and identify the most basic steps that the scripts are doing. I will be creating an example module to follow. Let's call it For example, if your scripts install a component, Some steps could be:

   1. Pre-requisite check.
   1. Installation of component.
   1. Validation of install.
   1. Uninstallation of component.

1. For each identified step, create a directory. For example, continuing with the above, your directory structure should look like this:

#### 3. (Optional) Use templating

#### 4. Build the binary

1. Clone the repo
1. `make install`
1. Make sure your `PATH` env variable has `go/bin` included (`export PATH=$PATH:$(go env GOPATH)/bin)`)

## CLI

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

## UI

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

### Error stack

- https://github.com/gruntwork-io/terratest/blob/master/modules/helm/template.go#L83
- https://godoc.org/github.com/go-errors/errors
- https://github.com/gruntwork-io/gruntwork-cli/blob/master/errors/errors.go

### Embedding static content

- https://github.com/rakyll/statik/pull/101/files
- https://github.com/golang/go/issues/41191
