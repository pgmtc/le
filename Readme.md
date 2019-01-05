# orchard-cli
[![Build Status](https://travis-ci.com/mara-mfa/pru-gh-search.svg?branch=master)](https://travis-ci.com/mara-mfa/pru-gh-search)
[![codecov](https://codecov.io/gh/pgmtc/orchard-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/pgmtc/orchard-cli)

## Installation
1. Download appropriate package from [releases page](https://github.com/pgmtc/orchard-cli/releases)
2. Unzip, there should be an executable
3. Put it somewhere to path
4. Run it from the terminal / command line

## Usage 
Syntax:

`orchard [module] [action] parameters`

In cases related to containers (vast majority), syntax is as follows:

- `orchard [module] [action] component` : runs action for component
- `orchard [module] [action] component1 component2 ... componentN` : runs for component1 .. componentN
- `orchard [module] [action] all` : runs for all available components


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
