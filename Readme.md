# orchard-cli
[![CircleCI](https://circleci.com/gh/pgmtc/orchard-cli.svg?style=svg)](https://circleci.com/gh/pgmtc/orchard-cli)

## Installation
1. TODO

## Usage 
Syntax:

`orchard [module] [action] parameters`

### Modules
#### local
Local module is responsible for running local environments
It has the following actions

`orchard local status`: prints status of the local environment

`orchard local pull [component]`: used for components with remote docker images

`orchard local create [component]`: create a docker container for the component

`orchard local remove [component]`: removes docker container of the component

`orchard local start [component]`: starts docker container for the component

`orchard local stop [component]`: stops the docker container for the component

`orchard local logs [component]`: shows logs of the related docker container

`orchard local watch [component]`: shows logs on the 'follow' basis

#### builder
Builder module is used for building containers. It has to be run from orchard-poc-umbrella directory

`orchard builder build [component]`: builds a docker image for the component


#### config
Config is a centralized storage used by other modules.

`orchard config init`: Run after the installation. Creates ~/.orchard, config file and default profile

`orchard config status'`: Prints out information about the current profile. Adding -v makes it more verbose

`orchard config create [profile] [source-profile]`: Creates a new profile. By passing source-profile parameter (not mandatory), it uses it as a base for copy

`orchard config switch [profile]`: Switches current profile to another one

### Modules TODO

#### source
Source module will be responsible for source code manipulation in orchard-poc-umbrella project
