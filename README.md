# Mozart

An orchestrator for scripts. (Python and Bash scripts supported for now).

Mozart is a simple utility to convert your independent, messy bunch of scripts into a well defined, orchestrated program!

**Note:** No Golang code changes required, it is a simple drop-in type utility for your scripts.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [What exactly does an orchestrator do?](#what-exactly-does-an-orchestrator-do)
  - [Life without an orchestrator](#life-without-an-orchestrator)
  - [Benefits of using Mozart](#benefits-of-using-mozart)
- [Using your scripts with Mozart](#using-your-scripts-with-mozart)
  - [Steps](#steps)
    - [1. Initial setup](#1-initial-setup)
    - [2. Modularize](#2-modularize)
    - [3. (Optional) Use templating](#3-optional-use-templating)
    - [4. Build the binary](#4-build-the-binary)
- [CLI](#cli)
  - [Mozart commands](#mozart-commands)
  - [Executing modules](#executing-modules)
  - [Checking the state](#checking-the-state)
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

Simply speaking, an orchestrator manages execution of all of your scripts.

### Life without an orchestrator

Suppose you had 2 bash scripts which let you install and uninstall a particular component respectively. Without an orchestrator:

1. You would have to ship both files to anyone who wants to use them. (imagine if you had 10 scripts, you would have to tar them up and send all).
1. If running on a shared system, you can never know if another person already ran the install script (unless actually going through the effort of seeing if the component was actually installed).
1. Tendency to use lesser scripts which makes each script huge. Breaking them down into multiple smaller scripts makes it easy to manage the code but it becomes harder to manage execution of all scripts.
1. No way of using common variables across the scripts (like version of the component being managed).
1. No way for one developer to look at the logs of a script executed by a second developer (unless the second developer explicitly routes the logs to an external file).
1. No way of accidentally preventing execution of a script more than once (like prevent install script from running again if it already ran once).
1. Manually execute each script through bash or python (no CLI or UI)

### Benefits of using Mozart

1. Simple migration to Mozart - `no Go code change necessary`. Just create a directory and dump your scripts in that - it's that easy. (Discussed in detail below)
1. `Single binary file` which contains all your scripts AND brings with it the CLI along with the UI - so no need to send in bunch of scripts to anyone anymore.
1. Lets you `modularize` the scripts, which means you can have more number of smaller scripts which do smaller tasks. No need to maintain huge bash files anymore. The smaller the scripts - the easier it is for you to manage and maintain them.
1. Ability to use `templating` capabilities - similar to helm. Values in the yaml file are accessible by all scripts.
1. `Public visibility of the state` of execution of scripts, so everyone has a clear idea of whether scripts were executed or not.
1. `Central logs` for everyone to see.
1. Make sure scripts are `executed just once` - the orchestrator will not allow you to run the same script again without explicitely mentioning a `Re-Run` flag, thereby preventing accidental execution of the same script.
1. Ready to use `CLI + UI` to manage the execution.

## Using your scripts with Mozart

Let us walk through how to actually add your scripts. There is one term that is of prime importance to Mozart - Modules.

**Modules**

Mozart works with the concept of modules and not the scripts directly. Modules are nothing but directories, created with the intention of performing 1 simple task. Each module (aka directory) can consist of either scripts or more nested modules (nested directories).

So in short, all your scripts need to be in a module, and Mozart will help you control the execution of those modules (instead of the scripts themselves).

**Sample module**

There is already a sample module present called `test-module` under `resources/templates`, which you can use to reference.

### Steps

#### 1. Initial setup

1. Clone the repo.
1. Create a new directory inside `resources/templates`. This will be the base module under which all your modules will exist.

#### 2. Modularize

1. Look through your scripts, and identify the most basic steps that the scripts are supposed to be performing. Suppose I want to install a component called `Symphony`. I have one huge bash script for that. Some basic steps inside that bash script could be:

   1. Pre-requisite check.
   1. Installation of component.
   1. Validation of install.
   1. Uninstallation of component.

1. For each identified step, create a module(directory) within the base directory and add that part of the script within that directory. For example, continuing with the above example, the directory structure should look like this:

```
resources/templates
├── symphony                    (this is the base module)
│   ├── 00-pre-req              (first sub-module)
│   │   └── pre-req.sh
│   ├── 10-install              (second sub-module)
│   │   ├── 00-install-step1.sh
│   │   └── 10-install-step2.sh
│   ├── 20-validate             (third sub-module)
│   │   └── validate.sh
│   └── 30-uninstall            (fourth sub-module)
│       └── uninstall.sh
```

The huge bash script is now broken down into smaller scripts, each in its own module. This makes the script easy to manage, while giving the option to add more scripts in the future as needed.

**Note:** the `xx-`prefix before a module or script name is an optional prefix, through this you can control the order of execution of scripts/modules within the module.

#### 3. (Optional) Use templating

You can do something like this within any script:

```
echo "{{.values.value1}} {{.values.value2}}"
```

These values are going to be fetched from a yaml file that you supply while invoking the CLI or the UI. (discussed later)

The yaml file will look something like this for the example above:

```
$ cat mozart-sample.yaml
values:
  value1: hello
  value2: world
```

When you execute the corresponding module that contains the above script, you will see

```
echo hello world
```

#### 4. Build the binary

`make build` -> for Darwin (Mac)

`make build-linux` -> for Linux

Voila! You have a single orchestrator binary with all your scripts in it.

## CLI

### Mozart commands

```
- execute       (executes all scripts in specified directory)
- state         (displays install state of all components, accepts optional args)
- version       (displays version info for the application)
  --json                (gives output in JSON)

Global flags
  -c                  (configuration file, defaults to 'mozart-sample.yaml')
  -v                  (prints verbosely, useful for debugging)
  -d, --dry-run       (optional) shows what scripts will run, but does not run the scripts
  -n, --no-generate   (optional) do not generate bash scripts as part of execution, instead use the ones in generated folder. Useful for running local change to the scripts
  -p, --parallel      (optional) Run all scripts in parallel
  -r, --re-run        (optional) re-run script from initial state, ignoring previously saved state

```

### Executing modules

Running the binary built in the earlier step, you will see something like this:

```
$ ./bin/mozart execute

*****************************************
Available commands:

mozart execute symphony-module
mozart execute symphony-module pre-req
mozart execute symphony-module install
mozart execute symphony-module validate
mozart execute symphony-module uninstall
*****************************************
```

If you select a module that contains other modules, something like `mozart execute symphony-module`, that's where the ordering of the sub-modules comes into play, which you control by adding the prefix. Or you can choose to execute a sub-module directly.

**Note**: Mozart automatically removes any prefix of the form `xx-`before the module name.

### Checking the state

Once you start the execution, the `state` command shows you the current state of execution of the various modules within Mozart, along with other information.

```
$ ./bin/mozart state

State: {
  "generated/symphony-module/00-pre-req": {
    "pre-req.sh": {
      "startTime": "2020-11-04T15:08:12.5981-08:00",
      "timeTaken": "8.218745ms",
      "lastSuccessTime": "2020-11-04 15:08:12.606306 -0800 PST m=+0.016747847",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.597-pre-req.log"
    }
  },
  "generated/symphony-module/10-install": {
    "00-install-step1.sh": {
      "startTime": "2020-11-04T15:08:12.607053-08:00",
      "timeTaken": "9.723926ms",
      "lastSuccessTime": "2020-11-04 15:08:12.616754 -0800 PST m=+0.027195761",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.606-00-install-step1.log"
    },
    "10-install-step2.sh": {
      "startTime": "2020-11-04T15:08:12.617443-08:00",
      "timeTaken": "7.333338ms",
      "lastSuccessTime": "2020-11-04 15:08:12.624767 -0800 PST m=+0.035208922",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.617-10-install-step2.log"
    }
  },
  "generated/symphony-module/20-validate": {
    "validate.sh": {
      "startTime": "2020-11-04T15:08:12.625411-08:00",
      "timeTaken": "7.542653ms",
      "lastSuccessTime": "2020-11-04 15:08:12.632945 -0800 PST m=+0.043386468",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.625-validate.log"
    }
  },
  "generated/symphony-module/30-uninstall": {
    "uninstall.sh": {
      "startTime": "2020-11-04T15:08:12.633649-08:00",
      "timeTaken": "8.040003ms",
      "lastSuccessTime": "2020-11-04 15:08:12.641679 -0800 PST m=+0.052120673",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.633-uninstall.log"
    }
  }
}
```

## UI

_Under development by Tosha Kamath_.

A sample UI wireframe:

![](https://user-images.githubusercontent.com/36335212/97355431-214dd300-1854-11eb-9103-0a39fcf3cf87.png)

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
