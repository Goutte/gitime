gitime : time tracker using git commit message commands
=======================================================

[![MIT](https://img.shields.io/github/license/Goutte/gitime?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Goutte/gitime?include_prereleases&style=for-the-badge)](https://github.com/Goutte/gitime/releases)
<!--
[![Build Status](https://img.shields.io/github/actions/workflow/status/Goutte/gitime/go.yml?style=for-the-badge)](https://github.com/Goutte/gitime/actions/workflows/go.yml)
[![Coverage](https://img.shields.io/codecov/c/github/Goutte/gitime?style=for-the-badge&token=FEUB64HRNM)](https://app.codecov.io/gh/Goutte/gitime/)
[![A+](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/Goutte/gitime)
-->
[![Code Quality](https://img.shields.io/codefactor/grade/github/Goutte/gitime?style=for-the-badge)](https://www.codefactor.io/repository/github/Goutte/gitime)


Purpose
-------

Collect, addition and return all the `/spend` and `/spent` time-tracking directives in git commit messages.

> This only looks at the `git log` of the currently checked out branch.


Download
--------

You can download the [latest build from the releases](https://github.com/Goutte/gitime/releases).

You can also install via `go get` (hopefully) :

```
go get -u github.com/goutte/gitime
```

or `go install`:

```
go install github.com/goutte/gitime
```

> If that fails, you can install by cloning and running `make install`.


Usage
-----

Go into your git-versioned project's directory, and run:

```
gitime sum
```

You can also get the spent time in a specific unit :

```
gitime sum --minutes
gitime sum --hours
```


Develop
-------

```
git clone https://github.com/Goutte/gitime.git
cd gitime
go get
go run main.go
```


Build & Run & Install
---------------------

```
make
make sum
make install
```

