![cloudlife](./images/cloudlife-logo-slogan.png)

What if [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) was a cloud-native serverless application in [WebAssembly](https://webassembly.org/) with [WASI](https://github.com/WebAssembly/WASI) using [TinyGo](https://tinygo.org/) and deployed using [Fermyon Spin](https://github.com/fermyon/spin)?

Welcome to cloudlife... "Life As A Service"

Uses [Vita](https://github.com/acifani/vita) for the Go language Game of Life implementation.

## Installation

- [Install Go 1.22+](https://go.dev/dl/)

- [Install TinyGo 0.35.0+](https://tinygo.org/getting-started/install/)

- [Install Fermyon Spin 3.0+](https://developer.fermyon.com/spin/v3/quickstart)

- Clone this repo

## Building the application

```
spin build
```

## Running the application

```
spin up
```

## Endpoints

### Multiverse

All operations on the multiverse are applied to all existing universes.

- `POST /multiverse?n=2`

    Creates new universes in a grid. Pass the `n` query parameter to create a specific number of universes.

    Returns the list of IDs for each universe that is created.

- `GET /multiverse`

    Returns the list of IDs for each universe that has been created.

- `PUT /multiverse`

    Advances each universe by 1 generation.

    Returns the cells for each universe after it has advanced.

- `DELETE /multiverse`

    Deletes all universes.

### Universe

All operations on a universe only apply to that specific universe.

- `POST /universe`

    Creates a new universe.

    Returns the ID for the new universe that is created.

- `GET /universe/:id`

    Returns the current cells for the universe with the specific ID.

- `PUT /universe/:id`

    Advances the universe with the specific ID by 1 generation.

    Returns the cells for the universe after it has advanced.

- `PUT /universe/:id?topid=:tid&bottomid=:bid&leftid=:leftid&rightid=:rid`

    Sets the neighbors for the universe with the specific ID to the universes with the respective ID and position e.g. top, bottom, left, and right. Does not advance the universe to the next generation.

- `DELETE /universe/:id`

    Deletes the universe with the specific ID.

## lifectl

`lifectl` is a command life tool to control cloudlife applications.

### Building

```
cd ./cmd/lifectl
go install .
```

### Running

```
$ lifectl
NAME:
   lifectl - Control your cloudlife

USAGE:
   lifectl [global options] command [arguments]

COMMANDS:
   start    Starts a cloudlife application
   run      Runs the cloudlife application
   stop     Stops a cloudlife application
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value  Host to use to connect to the cloudlife application (default: "http://localhost:3000")
   --help, -h    show help
```
