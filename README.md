![Tests](https://github.com/trntv/wilson/workflows/Test/badge.svg)
![Licence](https://img.shields.io/github/license/trntv/wilson)
[![Requirements Status](https://requires.io/github/trntv/wilson/requirements.svg?branch=master)](https://requires.io/github/trntv/wilson/requirements/?branch=master)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/trntv/wilson)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/trntv/wilson)
![GitHub closed issues](https://img.shields.io/github/issues-closed/trntv/wilson)
![GitHub issues](https://img.shields.io/github/issues/trntv/wilson)
![GitHub top language](https://img.shields.io/github/languages/top/trntv/wilson)
[![Go Report Card](https://goreportcard.com/badge/github.com/trntv/wilson)](https://goreportcard.com/report/github.com/trntv/wilson)

# Wilson - routine tasks automation toolkit
Wilson allows you to get rid of a bunch of bash scripts and to design you workflow pipelines in nice and neat way 
in yaml files. With Wilson you can design pipelines which are directed acyclic graphs that are made up of tasks and their dependencies on each other.

Automation is based on four concepts:
1. Execution context
2. Task
3. Pipeline that describes set of tasks to run
4. Optional watcher that listens for filesystem events and trigger tasks


## Install
### MacOS
```
brew install trntv/wilson/wilson
```
### Linux
```
curl -L https://github.com/trntv/wilson/releases/latest/download/wilson-linux-amd64.tar.gz | tar xz
```
### From sources
```
go get -u github.com/trntv/wilson
```

## Examples
### Tasks
[Task config](https://github.com/trntv/wilson/blob/master/example/task.yaml)
```
wilson -c example/task.yaml run task echo-date-local
wilson -c example/task.yaml run task echo-date-docker
``` 
### Pipelines
[Pipelines config](https://github.com/trntv/wilson/blob/master/example/pipeline.yaml)
```
wilson -c example/pipeline.yaml run test-pipeline
wilson -c example/pipeline.yaml run pipeline1
```

### Contexts
[Contexts config](https://github.com/trntv/wilson/blob/master/example/contexts.yaml)

### Watchers
[Watchers config](https://github.com/trntv/wilson/blob/master/example/watchers.yaml)
```
wilson -c watch.yaml --debug watch test-watcher test-watcher-2
```

### Full config
[Full config example](https://github.com/trntv/wilson/blob/master/example/full.yaml)

## Contexts
Available context types:
- local - shell
- container - docker, docker-compose, kubectl
- remote - ssh

## Pipelines
This configuration:
```yaml
pipelines:
    pipeline1:
        - task: start task
        - task: task A
          depends_on: "start task"
        - task: task B
          depends_on: "start task"
        - task: task C
          depends_on: "start task"
        - task: task D
          depends_on: "task C"
        - task: task E
          depends_on: ["task A", "task B", "task D"]
        - task: finish
          depends_on: ["task A", "task B", "finish"]
          
tasks:
    start task: ...
    task A: ...
    task B: ...
    task C: ...
    task D: ...
    task E: ...
    finish: ...
    
```
will create this pipeline:
```
               |‾‾‾ task A ‾‾‾‾‾‾‾‾‾‾‾‾‾‾|
start task --- |--- task B --------------|--- task E --- finish
               |___ task C ___ task D ___|
```

## Watchers
WIF*

## Autocomplete
### Bash
Add to  ~/.bashrc or ~/.profile
``
. <(wilson completion bash)
``

### ZSH
Add to  ~/.zshrc
``
. <(wilson completion zsh)
``

## Why "Wilson"?
https://en.wikipedia.org/wiki/Cast_Away#Wilson_the_volleyball 🏐

---
*waiting for inspiration
