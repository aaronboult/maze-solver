# maze-solver
A maze solver written in go implementing Dijkstra's algorithm. The mazes provided must be .png files.

### Image Requirements

All images must be .png images and must only contain pure black or pure white pixels, must contain a pure black outer border and may only contain two breaks in the outer border.

### Program Usage

Compile with `go build` in the *MazeSolver/* directory

### Usage:
`MazeSolver.exe [-maze] [-output] [options]`
### Options:
```
-showNumericalSolution: Show the path taken represented by a series of numbers
-preventWritingOutput: Prevents an output maze solution image being written
```
