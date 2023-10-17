# gh-wrun

This is an extension for [gh CLI](https://github.com/cli/cli) that allows manual interactive execution of workflows registered in GitHub Actions.

Interactive selection of workflows inputs in various formats.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![demo](https://github.com/t4kamura/gh-wrun/assets/51415522/c1c54584-fa08-45d7-907a-0bfeebf75fd1)

> **Note**
> Something similar can be done with the `gh workflow run` command,
> but there is no choice type select (at the moment).
> I thought about contributing to the `gh cli`,
> but since I was already running this tool personally, I decided to make it public.

## Installation

1. Install the `gh` CLI. [See how to install](https://github.com/cli/cli#installation)

   The `gh` cli must be version `2.35.0` or higher.
  (This is because the workflow list is obtained in json format.)

2. Install this extension

   ```sh
   gh extension install t4kamura/gh-wrun
   ```

## Usage

To get started

```sh
gh wrun
```

Execute this command in the root directory of the repository you wish to run.

> **Note**
> Manual execution may need to be enabled on the GitHub side if this is your first time doing it manually.

`--help` for other options.

## Todo

- [ ] Add loading when executing gh commands internally.
- [ ] Add a mode to wait for workflows to finish.
- [ ] Inputs supports the environments workflow.

## License

MIT
