# Mozart

If you struggle with maintaining a lot of scripts:

```
╰$ tree scripts
scripts
└── my-service
    ├── install
    │   ├── install-component1.sh
    │   ├── install-component2.sh
    │   └── install-component3.sh
    └── test
        ├── test-component1.sh
        ├── test-component2.sh
        └── test-component3.sh
```

Mozart converts all to a single binary, along with a CLI and and UI to manage them, with ZERO lines of code!

![](https://github.com/countertenor/mozart/blob/master/gh%20images/execution.jpg)

**CLI**:

```
Available CLI commands:

mozart execute my-service
mozart execute my-service install
mozart execute my-service install install-component1
mozart execute my-service install install-component2
mozart execute my-service install install-component3
mozart execute my-service test
mozart execute my-service test test-component1
mozart execute my-service test test-component2
mozart execute my-service test test-component3
```

## What is mozart?

Mozart is a simple drop-in (no go-coding required) utility to orchestrate and attach a CLI and a UI to your scripts, making your independent, messy bunch of scripts into a well defined, orchestrated program complete with a CLI and UI!

**Within minutes**, instead of having hundreds of different scripts, you can have a **single binary**, complete with a CLI and a UI, which includes all those scripts. All you need to do is to drop the scripts into a folder structure (explained below). That's it!

**Note:** No code changes required, it is a simple drop-in type utility for your scripts.

## What mozart is not

Mozart is **NOT** a replacement for tools like Ansible or Chef, it is not that mature (yet). Instead, think of Mozart as a simple utility to manage a bunch of scripts.

If you have a bunch of scripts lying about which you use to test out a particular component - think of Mozart.

If you have a bunch of scripts which let you deploy a particular program on some remote server - think of Mozart.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [What is mozart?](#what-is-mozart)
- [What mozart is not](#what-mozart-is-not)
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
- [Mozart yaml file](#mozart-yaml-file)
  - [Templating](#templating)
  - [Using common snippets across scripts](#using-common-snippets-across-scripts)
  - [Optional configuration parameters](#optional-configuration-parameters)
    - [Log sub-directory](#log-sub-directory)
    - [Exec source](#exec-source)
    - [Delims](#delims)
- [CLI](#cli)
  - [Mozart commands](#mozart-commands)
  - [Executing modules](#executing-modules)
  - [Checking the state](#checking-the-state)
    - [State of all modules](#state-of-all-modules)
    - [State of specific module](#state-of-specific-module)
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

Suppose you had 2 bash scripts which you use to test 2 different components. Without an orchestrator:

1. You would have to ship both files to anyone who wants to use them. (imagine if you had 10 scripts, you would have to tar them up and send all).
1. If running on a shared system, you can never know if another person already ran the scripts (unless actually going through the effort of seeing if the tests ran).
1. Tendency to use lesser number of scripts for easy script execution management, but that makes each script huge. Breaking them down into multiple smaller scripts makes it easy to manage the code but it becomes harder to manage execution of all scripts.
1. No way of using common variables across the scripts (like version of the component being tested).
1. No way for one developer to look at the logs of a script executed by a second developer (unless the second developer explicitly routes the logs to an external file).
1. No way of accidentally preventing execution of a script more than once (like prevent first script from running again if it already ran once).
1. Manually execute each script through bash or python (no CLI or UI)

### Benefits of using Mozart

1. Simple migration to Mozart - `no Go code change necessary`. Just create a directory and dump your scripts in that - it's that easy. (Discussed in detail below)
1. `Single binary file which contains all your scripts` AND brings with it the CLI along with the UI - so no need to send in bunch of scripts to anyone anymore.
1. Lets you `modularize` the scripts, which means you can have more number of smaller scripts which do smaller tasks. No need to maintain huge bash files anymore. The smaller the scripts - the easier it is for you to manage and maintain them.
1. Ability to use `templating` capabilities - similar to helm. Values in a `yaml` file are accessible by all scripts.
1. `Public visibility of the state` of execution of scripts, so everyone has a clear idea of whether scripts were executed or not.
1. `Central logs` for everyone to see.
1. Make sure scripts are `executed just once` - the orchestrator will not allow you to run the same script again without explicitly mentioning a `Re-Run` flag, thereby preventing accidental execution of the same script.
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

1. Install go - https://golang.org/doc/install
1. Clone repo - https://github.com/countertenor/mozart
1. Run the command `export PATH=$PATH:$(go env GOPATH)/bin`) (_Please note_ - to make it persistent, you will have to add the command to your .bashrc)
1. Run `go get github.com/rakyll/statik`
1. Create a new directory inside `resources/templates`. This will be the base module under which all your modules will exist.
1. (Optional) Delete the existing `test-module` inside `resources/templates` if you want a clean slate. That folder is only for reference. Leaving that folder as it is will not do any harm.

#### 2. Modularize

1.  Look through your scripts, and identify the most basic steps that the scripts are supposed to be performing. Suppose I want to install a component called `Symphony`. I have one huge bash script for that. Some basic steps inside that bash script could be:

    1. Pre-requisite check.
    1. Installation of component.
    1. Validation of install.
    1. Uninstallation of component.

1.  For each identified step, create a module(directory) within the base directory and add that part of the script within that directory. For example, continuing with the above example, the directory structure should look like this:

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

    The huge bash script is now broken down into smaller scripts, each in its own module. This makes the script easy to manage, while giving the option to add more scripts in the future as needed.

    **Note:** the `xx-`prefix before a module or script name is an optional prefix, through this you can control the order of execution of scripts/modules within the module.

**Note:** In case you don't want to break down your script into smaller scripts, you can create only the base module and drop your script in that directory.

#### 3. (Optional) Use templating

You can do something like this within any script:

    echo "{{.values.value1}} {{.values.value2}}"

These values are going to be fetched from a `yaml` file that you supply while invoking the CLI or the UI. (discussed later)

The `yaml` file should have something like this for the example above:

    values:
    value1: hello
    value2: world

When you execute the corresponding module that contains the above script, you will see

    echo hello world

This is discussed in detail below.

#### 4. Build the binary

Run `make` - builds binaries for linux, mac and centOS environments, inside bin directory.

Voila! You have a single orchestrator binary with all your scripts in it.

## Mozart yaml file

Providing an optional yaml file at runtime lets you enable certain templating features and configuration changes. If you do not need any changes, you can skip this section.

A sample blank configuration file can be generated from the binary itself, using the `init` command.

    ./bin/mozart init

    Generated sample file :  mozart-sample.yaml

### Templating

You can use the above `yaml` file to help with templating. If you refer to the `test-module` under `resources/templates`, you can see some examples of templating.

This idea is similar to helm.

**Example:**

In the file `step1.sh`, you see this:

    #!/bin/bash

    echo "{{.values.value1}} {{.values.value2}}"

The values in the brackets are the values that will be fetched from the `yaml` at runtime. So if you want to substitute some values at runtime, you replace the values with the `{{ }}` notation as you see above, and in the `yaml` file, add:

    values:
    value1: hello
    value2: world

Mozart will substitute these values at runtime.

### Using common snippets across scripts

There might be a scenario in which some scripts have a lot of common code. It is never a good idea to duplicate logic across scripts (DRY principle).

To tackle this, you can make use of the `common.yaml` file that is present in the `resources` folder. This file has one purpose and one purpose only - to hold common snippets of information that will be needed by more than one script.

You can take a look at the `resources/templates/test-module/10-python-module/00-module1/python-1.py` file for an example.

**Example:**

Suppose if you have a function that you want in more than one script, say

    def my_func(str):
    print(f'inside funct {str}')

Instead of having this function be duplicated across scripts, you add this function in the `common.yaml` file:

    my_func: >
    def my_func(str):
        print(f'inside funct {str}')

The `key` is `my_func`, and the value is the function itself.

You can then access this function in any script, using

    {{.my_func}}

**Note 1:** Sometimes you might want to add indentation to the above substituted lines of code (It is essential in python scripts). You can do so by using `nindent` (courtesy of [sprig functions](http://masterminds.github.io/sprig/strings.html))

    {{.my_func | nindent 4}}

**Note 2:** The only difference between the `common.yaml` and the main `yaml` file for Mozart config is that the `common.yaml` is more for compile time deduplication, whereas the main `yaml` file is for runtime changes. For example, functions that are duplicated will never need to be changed at runtime (common.yaml), whereas username and password should never be saved at compile time, instead should be provided at runtime.

### Optional configuration parameters

There are certain configuration parameters also that you can change using the same `yaml` file as above.

#### Log sub-directory

By default, log files are stored in `/var/log/mozart` directory (For linux and centOS), but if for some reason you want to add a sub-directory, you can do so by adding one line to the `yaml` file:

**Example:**

    log_path: my-log-dir

Then all the logs will go to:

    var/log/mozart/my-log-dir

#### Exec source

This lets you choose the execution environment of any type of script that you include.

The format is `file_ext: source`

**Example:**

    exec_source:
    py: /usr/bin/python
    sh: /bin/bash

This lets Mozart know that if you place any file with the extension of `.sh`, then run it using `/bin/bash`. If you place any file with the extension `.py`, then run it using `/usr/bin/python`.

**Note:** These are the only 2 extensions added by default in Mozart. If you add any other type of script apart from python or bash, you will need to add the execution source in the `yaml`.

#### Delims

This lets you change the default delimiters (default - `{{`, `}}`)

**Example 1:**

    delims: ["[[", "]]"]

Adding this line in the `yaml` file changes the delimiters to `[[ ]]`. So after this, you can use templating like:

    echo "[[.values.value1]] [[.values.value2]]"

**Example 2:**

    delims: ["<<", ">>"]

Adding this line in the `yaml` file changes the delimiters to `<< >> `. So after this, you can use templating like:

    echo "<<.values.value1>> <<.values.value2>>"

## CLI

Once you build your binary, Mozart gives you a CLI:

### Mozart commands

    - init          Generate a blank sample config yaml file for the orchestrator
    - execute       (executes all scripts in specified directory)
    - state         (displays install state of all modules, accepts optional args [module name])
    - server        Starts the REST server
    - version       (displays version info for the application)
    --json                (gives output in JSON)

    Global flags
    -c                  (optional) (configuration file, defaults to 'mozart-sample.yaml')
    -v                  (optional) (prints verbosely, useful for debugging)
    -d, --dry-run       (optional) shows what scripts will run, but does not run the scripts
    -n, --no-generate   (optional) do not generate bash scripts as part of execution, instead use the ones in generated folder. Useful for running local change to the scripts
    -p, --parallel      (optional) Run all scripts in parallel
    -r, --re-run        (optional) re-run script from initial state, ignoring previously saved state

### Executing modules

Running the binary built in the earlier step, you will see something like this:

```
$ ./bin/mozart execute -h

Execute scripts inside any folder.
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

#### State of all modules

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

#### State of specific module

To get the state of a particular module:

```
./mozart state validate

State: {
  "generated/symphony-module/20-validate": {
    "validate.sh": {
      "startTime": "2020-11-04T15:08:12.625411-08:00",
      "timeTaken": "7.542653ms",
      "lastSuccessTime": "2020-11-04 15:08:12.632945 -0800 PST m=+0.043386468",
      "lastErrorTime": "",
      "state": "success",
      "logFilePath": "logs/2020-11-04--15-08-12.625-validate.log"
    }
  }
```

## UI

Developed by @toshakamath

![](https://github.com/countertenor/mozart/blob/master/gh%20images/1.jpg)
![](https://github.com/countertenor/mozart/blob/master/gh%20images/2.jpg)

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
- https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff

### Error stack

- https://github.com/gruntwork-io/terratest/blob/master/modules/helm/template.go#L83
- https://godoc.org/github.com/go-errors/errors
- https://github.com/gruntwork-io/gruntwork-cli/blob/master/errors/errors.go

### Embedding static content

- https://github.com/rakyll/statik/pull/101/files
- https://github.com/golang/go/issues/41191
- https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/
- https://github.com/akmittal/go-embed (example for react app)
