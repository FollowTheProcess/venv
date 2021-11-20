# venv

[![License](https://img.shields.io/github/license/FollowTheProcess/venv)](https://github.com/FollowTheProcess/venv)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/venv)](https://goreportcard.com/report/github.com/FollowTheProcess/venv)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/venv?logo=github&sort=semver)](https://github.com/FollowTheProcess/venv)
[![CI](https://github.com/FollowTheProcess/venv/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/venv/actions?query=workflow%3ACI)

CLI to take the pain out of python virtual environments 🛠

* Free software: Apache Software License 2.0

## Project Description

`venv` aims to take all the pain and hastle out of creating and managing python virtual environments, as well
as installing project dependencies when working with python packages!

It does this by trying to "figure out" what to do based on the project context and the directories and files around it. If it
can't figure it out, it will ask you what you want it to do next.

`venv` used to be a (very large) shell function and I thought it best to rewrite it in a "proper" language, so here it is :tada:

## Installation

There are pre-built binaries in the [GitHub Releases] section and a homebrew tap:

```shell
brew tap FollowTheProcess/homebrew-tap

brew install FollowTheProcess/homebrew-tap/venv
```

## Quickstart

There's only one thing to do with `venv`, and that's run it!

```shell
venv
```

## Logic

The logical flow thet `venv` goes through to determine what to do with your project is as follows:

1. First it will look to see if there is a `.venv` or a `venv` directory under the current working directory. If there is it will simply say so and exit (unlike in shell scripts, an external program cannot alter the state of the shell that launched it, so we can't activate it for you sorry!)
2. It will then look for a `requirements_dev.txt`, because in projects where this exists, it typically contains everything needed to work on it. That's why we prefer `requirements_dev.txt` over plain old `requirements.txt`. If it finds one, it will create a python virtual environment and install the requirements from the file.
3. Failing that, we repeat the same process just this time with the classic `requirements.txt`
4. Now it looks for a `pyproject.toml`, and will do a few different things if it finds one:
   1. If it finds a `pyproject.toml` with either a `setup.cfg` or a `setup.py`, it knows that the project is based on [setuptools] and will install the project as such. If the setuptools file is a `setup.cfg`, it will attempt to install with `[dev]` extras, falling back to a normal install in all other cases.
   2. If it finds a `pyproject.toml` on it's own, it checks whether or not the file specifies a [poetry] or a [flit] based project. Making the appropriate call to whichever it finds
5. Now we're out of ideas! If we get here, `venv` will announce it cannot auto-detect the appropriate environment and ask you what you want to do next! You'll have the option to create a new environment or simply exit and take manual control

All output from the underlying calls is exposed back to the terminal so you can see everything that is happening. If you want some additional debugging information, you can set the `VENV_DEBUG` environment variable to 1 before running the program:

```shell
VENV_DEBUG=1 venv
```

[GitHub Releases]: https://github.com/FollowTheProcess/venv/releases
[poetry]: https://python-poetry.org
[flit]: https://flit.readthedocs.io/en/latest/
[setuptools]: https://setuptools.pypa.io/en/latest/
