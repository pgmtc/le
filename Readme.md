# local environment (le)
###### master
| | | |
--- | --- | ---
master  | [![Build Status](https://travis-ci.com/pgmtc/le.svg?branch=master)](https://travis-ci.com/pgmtc/le) | [![codecov](https://codecov.io/gh/pgmtc/le/branch/master/graph/badge.svg)](https://codecov.io/gh/pgmtc/le) | 
develop | [![Build Status](https://travis-ci.com/pgmtc/le.svg?branch=develop)](https://travis-ci.com/pgmtc/le) | [![codecov](https://codecov.io/gh/pgmtc/le/branch/develop/graph/badge.svg)](https://codecov.io/gh/pgmtc/le) |

## Prerequisites
1. Running docker daemon
2. AWS cli

## Installation
### Semi-automatic (experimental to be used on macOS)
* Download 1.0.1 package and store in /tmp

`curl -L  "https://github.com/pgmtc/le/releases/download/1.2.0/le_1.2.0_macOS_i386.tar.gz" | gunzip -c | tar -C /tmp -xvf -`

* Run update which should self-update and install into /usr/local/bin. Ignore errors about missing config for now

`/tmp/le config update-cli`

* Init config

`le config init`

### Manual
1. Download appropriate package from [releases page](https://github.com/pgmtc/le/releases)
2. Unzip it, there should be an executable inside
3. Put it somewhere to path, potentially chmod a+x it
4. Run it from the terminal / command line


## Usage 
Syntax:

`le [module] [action] parameters`

In cases related to containers (vast majority), syntax is as follows:

- `le [module] [action] component` : runs action for component
- `le [module] [action] component1 component2 ... componentN` : runs for component1 .. componentN
- `le [module] [action] all` : runs for all available components


## Modules
### local
Local module is responsible for running local environments
It has the following actions

`le local status`: prints status of the local environment

`le local pull [component]`: used for components with remote docker images

`le local create [component]`: create a docker container for the component

`le local remove [component]`: removes docker container of the component

`le local start [component]`: starts docker container for the component

`le local stop [component]`: stops the docker container for the component

`le local logs [component]`: shows logs of the related docker container

`le local watch [component]`: shows logs on the 'follow' basis

### builder
`orchard builder build [component]`: builds a docker image for the component

Build definition is stored in .builder directory inside the project.
It can be override by providing --specdir argument.
To create empty build definition directory, run 
`le builder init`


### config
Config is a centralized storage used by other modules.

`le config init`: Run after the installation. Creates ~/.le, config file and default profile

`le config status'`: Prints out information about the current profile. Adding -v makes it more verbose

`le config create [profile] [source-profile]`: Creates a new profile. By passing source-profile parameter (not mandatory), it uses it as a base for copy

`le config switch [profile]`: Switches current profile to another one
