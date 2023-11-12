# cloudlife

What if [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) was a cloud-native serverless application written for [WASI](https://github.com/WebAssembly/WASI) using [TinyGo](https://tinygo.org/) and deployed using [Fermyon Spin](https://github.com/fermyon/spin)?

Uses [Vita](https://github.com/acifani/vita) for the Go language Game of Life implementation.

## Building

```
spin build
```

## Running

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
