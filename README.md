# Mozart

An orchestrator for scripts. (Python and Bash scripts supported for now).

Mozart is a simple utility to convert your independent, messy bunch of scripts into a well defined, orchestrated program!

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [What exactly does an orchestrator do?](#what-exactly-does-an-orchestrator-do)
- [Benefits of using Mozart](#benefits-of-using-mozart)
- [Easy steps to migrate to Mozart](#easy-steps-to-migrate-to-mozart)
- [Using your scripts with Mozart](#using-your-scripts-with-mozart)
  - [Steps](#steps)
    - [1. Initial setup](#1-initial-setup)
    - [2. Modularize](#2-modularize)
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

## What exactly does an orchestrator do?

Good question. Simply speaking, an orchestrator manages execution of all of your scripts.

Suppose you had 2 bash scripts which let you install and uninstall a particular component respectively. Without an orchestrator:

1. You would have to ship both files to anyone who wants to use them. (imagine if you had 10 scripts, you would have to tar them up and send all).
1. If running on a shared system, you can never know if another person already ran the install script (unless actually going through the effort of seeing if the component was actually installed).
1. No way of using common variables across the scripts (like version of the component being managed).
1. No way of accidentally preventing execution of a script more than once (like prevent install script from running again if it already ran once).
1. Manually execute each script through bash or python (no CLI or UI)

## Benefits of using Mozart

- Simple migration to Mozart - `no Go code change necessary`. Just create a directory and dump your scripts in that - it's that easy. (Discussed in detail below)
- `Single binary` which contains all your scripts - so no need to send in bunch of scripts to anyone anymore.
- Lets you `modularize` the scripts, which means you can have more number of smaller scripts which do smaller tasks. No need to maintain huge bash files anymore. The smaller the scripts - the easier it is for you to manage and maintain them.
- Ability to use `templating` capabilities - similar to helm. Values in the yaml file are accessible by all scripts.
- `Public visibility of the state` of execution of scripts, so everyone has a clear idea of whether scripts were executed or not.
- Make sure scripts are `executed just once` - the orchestrator will not allow you to run the same script again without explicitely mentioning a `Re-Run` flag, thereby preventing accidental execution of the same script.
- Ready to use `CLI + UI` to manage the execution.

## Using your scripts with Mozart

Let us walk through how to actually add your scripts. There are some terms that are going to be used:

**Modules**

Mozart works with the concept of modules and not the scripts directly. Modules are nothing but directories, created with the intention of performing 1 simple task. Each module (or directory) can consist of either scripts or more nested modules (directories).

So in short, all your scripts need to belong to a module, and Mozart will help you control the execution of those modules.

**Sample module**

There is already a sample module present there called `test-module`, which you can use to reference.

### Steps

#### 1. Initial setup

1. Clone the repo
1. Create a base directory inside `resources/templates`. This will be the main directory under which all your scripts will exist.

#### 2. Modularize

1. Look through your scripts, and identify the most basic steps that the scripts are supposed to be performing. Suppose I want to install a component called `Symphony`. I have one huge bash script that does everything. Some steps inside that huge bash scripot could be:

   1. Pre-requisite check.
   1. Installation of component.
   1. Validation of install.
   1. Uninstallation of component.

1. For each identified step, create a module(directory) within the base directory and add that part of the script within that directory. For example, continuing with the above, your directory structure should look like this:

```
resources/templates
├── symphony
│   ├── 00-pre-req
│   │   └── pre-req.sh
│   ├── 10-install
│   │   ├── 00-install-step1.sh
│   │   └── 10-install-step2.sh
│   ├── 20-validate
│   │   └── validate.sh
│   └── 30-uninstall
│       └── uninstall.sh
```

**Note:** the `xx-`prefix before a module or script name is an optional prefix, through this you can control the order of execution within your module.

#### 3. (Optional) Use templating

You can do something like this within any script:

```
echo "{{.values.value1}} {{.values.value2}}"
```

These values are going to be fetched from a yaml file that you supply while invoking the CLI or the UI.

The yaml file will look something like this for the example above:

```
values:
  value1: hello
  value2: world
```

When you execute the corresponding module that contains the above script, you will see

```
echo hello world
```

#### 4. Build the binary

`make build-linux` -> for Linux
`make build` -> for Darwin (Mac)

Voila! You have a single orchestrator binary with all your scripts in it!

### CLI

```
mozart commands

- execute       (executes all scripts in specified directory)
- state         (displays install state of all components, accepts optional args)
- version       (displays version info for the application)
  --json                (gives output in JSON)

Global flags
  -c           (configuration file, defaults to 'mozart-sample.yaml')
  -v           (prints verbosely, useful for debugging)
  -d, --dry-run       (optional) shows what scripts will run, but does not run the scripts
  -n, --no-generate   (optional) do not generate bash scripts as part of execution, instead use the ones in generated folder. Useful for running local change to the scripts
  -r, --re-run        (optional) re-run script from initial state, ignoring previously saved state

```

Running the binary built in the earlier step, you will see something like this:

```
./bin/mozart execute

*****************************************
Available commands:

mozart execute symphony-module
mozart execute symphony-module pre-req
mozart execute symphony-module install
mozart execute symphony-module validate
mozart execute symphony-module uninstall
*****************************************
```

If you select a module that contains other modules, something like `mozart execute symphony-module`, that's where the ordering of the sub-modules comes into play, which you control by adding the prefix.

**Note**: Mozart automatically removes any prefix of the form `xx-`before the module name.

### UI

_Coming soon_

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
