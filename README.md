# SUDOKU

## Puzzle Structure

A puzzle is stored as JSON. `0` represents an empty cell.

```json
{
  "board": [
    [0, 0, 6, 0, 0, 0, 0, 8, 7],
    [7, 0, 9, 0, 4, 0, 0, 0, 5],
    [1, 0, 0, 0, 0, 7, 0, 2, 0],
    [0, 6, 0, 0, 0, 0, 5, 0, 0],
    [2, 4, 0, 8, 0, 0, 0, 0, 1],
    [0, 0, 3, 0, 0, 0, 0, 0, 6],
    [0, 0, 0, 0, 0, 0, 0, 5, 0],
    [6, 0, 8, 9, 5, 2, 0, 0, 3],
    [0, 5, 2, 0, 0, 3, 0, 0, 0]
  ]
}
```

## Commands

### Generate

Generate a 9x9 sudoku puzzle with a given number of givens. Saves to `puzzle-{timestamp}.json`.

```sh
go run main.go generate --n 30
```

`--n` sets the number of given squares. Clamped to 17–81. Default is 30.

| `--n` | Difficulty |
|-------|------------|
| 30+ | Easy (fast generation) |
| 28 | Moderate |
| 25 | Hard |
| 17 | Minimum givens (may be slow) |

### Read

Display a saved puzzle from a JSON file in a human-readable grid.

```sh
go run main.go read puzzle-1782579042.json
```

Example output:

```
. . 6 | . . . | . 8 7
7 . 9 | . 4 . | . . 5
1 . . | . . 7 | . 2 .
------+-------+------
. 6 . | . . . | 5 . .
2 4 . | 8 . . | . . 1
. . 3 | . . . | . . 6
------+-------+------
. . . | . . . | . 5 .
6 . 8 | 9 5 2 | . . 3
. 5 2 | . . 3 | . . .
```
