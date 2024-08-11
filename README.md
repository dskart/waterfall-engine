![release](https://img.shields.io/github/v/release/dskart/waterfall-engine)

# :ocean: Waterfall Engine :ocean:

[Live Demo](https://waterfall-engine.raphaelvanhoffelen.com/)

A Waterfall is a method by which the profits from investments are distributed amoung the various stakeholders. This distribution structure specifies the order and rules for how returns are shared between limited partners (LPs), who are typically the investores providing capital, and the general partners (GPs), who manage the private equity fund.

This repository contains the code for a waterfall engine that can be used to calculate the distribution of profits in a waterfall structure. The engine is written in Go. HTMX is used for the front-end.

The actual waterfall engine logic is implemented in the [./app/engine](./app/engine) module. 

## Getting Started

### Pre-requisites

- [Go 1.22](https://go.dev/doc/install)
- [Node.js v20.9.0](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

You will then need to download the tools for the project. You can do this by running the following commands:

```bash
# this will download the tools inside the /bin folder
make bin
```

### Input Data

In order to run correctly, the applications needs to load a [transactions.csv](./data/transactions.csv) and a [commitments.csv](./data/commitments.csv) [data](/data/) folder. Default data is provided in the [data](/data/) folder of the project and is loaded by default when running the application.

You can also pass in a custom data folder path using the `--data` flag:

```bash
go run main.go serve --data /path/to/data
```

### Running from Source

Before running the application, you will need to create a `config.yml` in the root directory of the project. The `config.yml` file should contain the following:

```yaml
App:
  Store:
    InMemory: true
  Engine:
    PreferredReturn:
      HurdlePercentage: 0.08
    CatchUp:
      CatchupPercentage: 1.0
      CarriedInterestPercentage: 0.2
    FinalSplit:
      LpPercentage: 0.8
      GpPercentage: 0.2
```

You then need to build the UI dependencies:

```bash
# this will download htmx, generate templ files, build css and js files
make ui
```

You can run the server directly from go and serve the application on [http://localhost:8080](http://localhost:8080):

```bash
go run main.go serve
```

Or you can use hot reload with air and browser-sync and serve the application on [http://localhost:3001](http://localhost:3001):

```bash
make serve
```

Run the following command to get a list of all available commands:

```bash
go run main.go --help
```

### Running from Docker

You can also run the application using Docker. To do this simply run the following script:

```bash
./run_docker.sh
```

This script will automatically build the Docker image and run the application in a container using the default waterfall configuration and the [./data](./data) folder as a mounted data source.

## Contributing

Commits should follow the following convention: `refactor|feat|fix|docs|breaking|misc|chore|test: description`

## CI/CD

Github actions is used to run the CI/CD pipeline. Workflows are found under [.github/workflows](.github/workflows).

A docker image is built and pushed to the [github packages](https://github.com/dskart/waterfall-engine/pkgs/container/waterfall-engine) on every commit and release.

## Versioning

This project uses [SemVer](https://semver.org/) for versioning.
The project also uses [release-please](https://github.com/googleapis/release-please) to automatically create releases based on the commit messages.

The project follows the [agile flow](https://medium.com/@mike.naquin/agileflow-4034af0b59fd) method for git branching and releasing

See releases [here](https://github.com/dskart/waterfall-engine/releases)

