# maze-solver
A maze solver written in go implementing Dijkstra's algorithm. The mazes provided must be .png files.

## Image Requirements

All images must be .png images and must only contain pure black or pure white pixels, must contain a pure black outer border and may only contain two breaks in the outer border.

## Program Usage

Compile with `go build` in the MazeSolver/ directory

```
Usage:
    MazeSolver.exe [-maze maze.png] [-output maze_output.png] [-showNumericalSolution] [preventWritingOutput]
```
