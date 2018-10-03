# Ionize
[![Build Status](https://travis-ci.org/ion-channel/ionize.svg?branch=master)](https://travis-ci.org/ion-channel/ionize)
[![Go Report Card](https://goreportcard.com/badge/github.com/ion-channel/ionize)](https://goreportcard.com/report/github.com/ion-channel/ionize)
[![GoDoc](https://godoc.org/github.com/ion-channel/ionize?status.svg)](https://godoc.org/github.com/ion-channel/ionize)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/ion-channel/ionize/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/ion-channel/ionize.svg)](https://github.com/ion-channel/ionize/releases/latest)

Wrapper around ion-connect to manage the asynchronous calls necessary to embed in a CI/CD process

# Requirements
Golang Version 1.8 or higher

# Installation
Ionize can be run either as a native tool or indirectly within a Docker
container.

Go:
```
go get github.com/ion-channel/ionize
```


Docker:
```
docker pull ionchannel/ionize
```

# Running
Ionize requires a key to authenticate with the Ion Channel API.  You can create one inside
the Ion Channel console.  Once you have the key you can supply it to Ionize with an environment
variable.

Running with the native tool:
```
IONCHANNEL_SECRET_KEY=<secret> ionize help
```



And within a docker container:
```
docker run -it -e IONCHANNEL_SECRET_KEY=<secret> ionchannel/ionize help
```



In addition to the api key you will also need a `.ionize.yaml` file in the current
working directory. The file contains ids for the project in Ion Channel to analyze as
well as any configuration needed.  An example can be seen [here](https://github.com/ion-channel/ionize/blob/master/.ionize.yaml.example).

# Versioning

The project will be versioned in accordance with [Semver 2.0.0](http://semver.org).  See the [releases](https://github.com/ion-channel/ionic/releases) section for the latest version.  Until version 1.0.0 the project is considered to be unstable.

# License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).  See [LICENSE.md](./LICENSE.md) for more information.
