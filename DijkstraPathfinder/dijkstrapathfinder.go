package dijkstrapathfinder

/*
	A package to find a path through a maze when given a jagged integer slice
	representing the maze where '1' is a wall and '0' is a path
*/

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
)

// DijkstraPathfinder handles solving a given maze, the end points, paths and logging
type DijkstraPathfinder struct {
	Maze      [][]int // Stores the entire maze for local use
	Logging   bool
	xEnd      int
	yEnd      int
	pathStack []path
}

// Solve beings the solving process
func (pathfinder *DijkstraPathfinder) Solve() ([][]int, int) {

	borderBreaks := pathfinder.getBorderBreaks()

	if len(borderBreaks) != 2 {

		fmt.Println("The borders of the maze may only contain 2 break points")

		os.Exit(1)

	}

	pathfinder.xEnd = borderBreaks[1][0]

	pathfinder.yEnd = borderBreaks[1][1]

	pathfinder.pathStack = []path{
		path{
			length: 0,
			pathTo: []node{
				node{
					xPos: borderBreaks[0][0],
					yPos: borderBreaks[0][1],
				},
			},
		},
	}

	return pathfinder.next()

}

// Get the points at which the maze starts and ends
func (pathfinder *DijkstraPathfinder) getBorderBreaks() [][]int {

	borderBreaks := [][]int{}

	for x := 0; x < len(pathfinder.Maze[0])-1; x++ {

		if pathfinder.Maze[0][x] == 0 {

			borderBreaks = append(borderBreaks, []int{x, 0})

		}

		if pathfinder.Maze[len(pathfinder.Maze)-1][x] == 0 {

			borderBreaks = append(borderBreaks, []int{x, len(pathfinder.Maze) - 1})

		}

	}

	for y := 0; y < len(pathfinder.Maze)-1; y++ {

		if pathfinder.Maze[y][0] == 0 {

			borderBreaks = append(borderBreaks, []int{0, y})

		}

		if pathfinder.Maze[y][len(pathfinder.Maze[0])-1] == 0 {

			borderBreaks = append(borderBreaks, []int{len(pathfinder.Maze[0]) - 1, y})

		}

	}

	return borderBreaks

}

// Test the next generation of paths
func (pathfinder *DijkstraPathfinder) next() ([][]int, int) {

	if pathfinder.pathStack[0].length == -1 {

		fmt.Println("No solution could be found")

		os.Exit(1)

	}

	newPaths, endPath := pathfinder.pathStack[0].extend(pathfinder.Maze, pathfinder.xEnd, pathfinder.yEnd)

	if endPath.length == -2 {

		pathfinder.pathStack = append(pathfinder.pathStack, newPaths...)

		pathfinder.sortPathStack()

		if pathfinder.Logging {

			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()

			fmt.Println(pathfinder.pathStack)

			fmt.Print("\n\n")

		}

		return pathfinder.next()

	}

	return endPath.getPath(), len(pathfinder.pathStack)

}

// Sort the slice of paths to have the shortest path on the top - Optimal solution
func (pathfinder *DijkstraPathfinder) sortPathStack() {

	sort.Slice(pathfinder.pathStack, func(i, j int) bool {
		if pathfinder.pathStack[j].length == -1 {
			return true
		}
		if pathfinder.pathStack[i].length == -1 {
			return false
		}
		return pathfinder.pathStack[i].length < pathfinder.pathStack[j].length
	})

}

// Individual paths, originating from the 'start' position
type path struct {
	length int
	pathTo []node
}

// Extend the number of routes through the maze to include the nodes of the closest adjacent nodes
func (currentPath *path) extend(maze [][]int, xEnd int, yEnd int) ([]path, path) {

	if currentPath.length != -1 {

		adj := currentPath.pathTo[len(currentPath.pathTo)-1].getAdjacentNodes(maze, currentPath.pathTo, xEnd, yEnd)

		if len(adj) > 0 {

			var newPaths []path = []path{}

			for i := 1; i < len(adj); i++ {

				newPaths = append(newPaths, path{
					length: currentPath.length + currentPath.pathTo[len(currentPath.pathTo)-1].getLength(adj[i], xEnd, yEnd),
					pathTo: make([]node, len(currentPath.pathTo)),
				})

				copy(newPaths[len(newPaths)-1].pathTo, currentPath.pathTo)

				newPaths[len(newPaths)-1].pathTo = append(newPaths[len(newPaths)-1].pathTo, adj[i])

			}

			currentPath.pathTo = append(currentPath.pathTo, adj[0])

			currentPath.length += currentPath.pathTo[len(currentPath.pathTo)-1].getLength(adj[0], xEnd, yEnd)

			lastNodeInPath := currentPath.pathTo[len(currentPath.pathTo)-1]

			endPath := path{length: -2} // An of 'path' with a length of -2 simply denotes the maze has not been solved

			if lastNodeInPath.xPos == xEnd && lastNodeInPath.yPos == yEnd {

				endPath = *currentPath

			}

			if endPath.length == -2 {

				for i := 0; i < len(newPaths); i++ {

					lastNodeInPath = newPaths[i].pathTo[len(currentPath.pathTo)-1]

					if lastNodeInPath.xPos == xEnd && lastNodeInPath.yPos == yEnd {

						endPath = newPaths[i]

						break

					}

				}

			}

			return newPaths, endPath

		}

		currentPath.length = -1

	}

	return []path{}, path{length: -2} // An of 'path' with a length of -2 simply denotes the maze has not been solved

}

func (currentPath *path) getPath() [][]int {

	pathTo := [][]int{}

	for _, currentNode := range currentPath.pathTo {

		pathTo = append(pathTo, []int{currentNode.xPos, currentNode.yPos})

	}

	return pathTo

}

// Junction points; one or more turns can be made
type node struct {
	xPos int
	yPos int
}

// Search for nodes that can be reached on a straight path
func (currentNode *node) getAdjacentNodes(maze [][]int, currentPath []node, xEnd int, yEnd int) []node {

	directionsToTry := currentNode.getDirectionsToTry(len(maze)-1, currentPath)

	newNodes := []node{}

	for _, directionToMove := range directionsToTry {

		for i := 0; i < len(maze)-1; i++ {

			if currentNode.canMove(maze, directionToMove) {

				if currentNode.xPos+directionToMove.XDirection == xEnd && currentNode.yPos+directionToMove.YDirection == yEnd {

					newNodes = append(newNodes, node{
						xPos: currentNode.xPos + directionToMove.XDirection,
						yPos: currentNode.yPos + directionToMove.YDirection,
					})

					return newNodes

				}

				if currentNode.nodeIsPresent(maze, directionToMove) {

					newNodes = append(newNodes, node{
						xPos: currentNode.xPos + directionToMove.XDirection,
						yPos: currentNode.yPos + directionToMove.YDirection,
					})

					break

				}

			} else {

				break

			}

			directionToMove.Increment()

		}

	}

	return newNodes

}

// Look for a possible change in direction at the current position
func (currentNode *node) nodeIsPresent(maze [][]int, directionToMove Direction) bool {

	tempNode := node{
		xPos: currentNode.xPos + directionToMove.XDirection,
		yPos: currentNode.yPos + directionToMove.YDirection,
	}

	if directionToMove.XDirection != 0 {

		return tempNode.canMove(maze, Direction{
			XDirection: 0,
			YDirection: 1,
		}) || tempNode.canMove(maze, Direction{
			XDirection: 0,
			YDirection: -1,
		})

	}

	return tempNode.canMove(maze, Direction{
		XDirection: 1,
		YDirection: 0,
	}) || tempNode.canMove(maze, Direction{
		XDirection: -1,
		YDirection: 0,
	})

}

func (currentNode *node) canMove(maze [][]int, directionToMove Direction) bool {

	targetX := currentNode.xPos + directionToMove.XDirection

	targetY := currentNode.yPos + directionToMove.YDirection

	if targetY == -1 || targetX == -1 || targetY == len(maze) || targetX == len(maze[0]) {

		return false

	}

	return maze[targetY][targetX] == 0

}

func (currentNode *node) getDirectionsToTry(mazeLength int, currentPath []node) []Direction {

	directionsToTry := []Direction{}

	if len(currentPath) == 1 {

		directionsToTry = append(directionsToTry, Direction{
			XDirection: 0,
			YDirection: 0,
		})

		switch currentPath[0].xPos {
		case 0: // From left of maze
			directionsToTry[0].XDirection = 1
			break
		case mazeLength: // From right of maze
			directionsToTry[0].XDirection = -1
			break
		}

		switch currentPath[0].yPos {
		case 0: // From top of maze
			directionsToTry[0].YDirection = 1
			break
		case mazeLength: // From bottom of maze
			directionsToTry[0].YDirection = -1
			break
		}

	} else {

		previousNode := currentPath[len(currentPath)-2]

		if currentNode.xPos-previousNode.xPos != 0 {

			directionsToTry = append(directionsToTry, []Direction{
				Direction{0, 1},  // Down
				Direction{0, -1}, // Up
			}...)

			if currentNode.xPos-previousNode.xPos > 0 { // If true, don't move to the left

				directionsToTry = append(directionsToTry, []Direction{
					Direction{1, 0}, // Right
				}...)

			} else { // If true, don't move to the right

				directionsToTry = append(directionsToTry, []Direction{
					Direction{-1, 0}, // Left
				}...)

			}

		} else {

			directionsToTry = append(directionsToTry, []Direction{
				Direction{1, 0},  // Right
				Direction{-1, 0}, // Left
			}...)

			if currentNode.yPos-previousNode.yPos > 0 { // If true, don't move upwards

				directionsToTry = append(directionsToTry, []Direction{
					Direction{0, 1}, // Down
				}...)

			} else { // If true, don't move downwards

				directionsToTry = append(directionsToTry, []Direction{
					Direction{0, -1}, // Up
				}...)

			}

		}

	}

	return directionsToTry

}

// Get the distance between the current node an a given adjacent node
func (currentNode *node) getLength(adjNode node, xEnd int, yEnd int) int {
	// The absolute value of Dx + Dy
	return int(math.Abs(float64(adjNode.xPos-currentNode.xPos))+math.Abs(float64(adjNode.xPos-currentNode.xPos))) +
		int(math.Abs(float64(xEnd-adjNode.xPos))+math.Abs(float64(yEnd-adjNode.yPos)))

}

// Direction represents the X and Y values to be travelled in
type Direction struct {
	XDirection int
	YDirection int
}

// Increment moves the current direction further away from it's point of origin
func (currentDirection *Direction) Increment() {

	if currentDirection.XDirection == 0 {

		if currentDirection.YDirection > 0 {

			currentDirection.YDirection++

		} else {

			currentDirection.YDirection--

		}

	} else {

		if currentDirection.XDirection > 0 {

			currentDirection.XDirection++

		} else {

			currentDirection.XDirection--

		}

	}

}

// Decrement moves the current direction towards being X: 0, Y: 0
func (currentDirection *Direction) Decrement() bool {

	if currentDirection.XDirection == 0 {

		if currentDirection.YDirection > 0 {

			currentDirection.YDirection--

		} else {

			currentDirection.YDirection++

		}

	} else {

		if currentDirection.XDirection > 0 {

			currentDirection.XDirection--

		} else {

			currentDirection.XDirection++

		}

	}

	return !(currentDirection.XDirection == 0 && currentDirection.YDirection == 0)

}
