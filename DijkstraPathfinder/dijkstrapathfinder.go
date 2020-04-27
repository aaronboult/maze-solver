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

/*
	The main maze solver type
	Uses Dijkstra's Algorithm with A*
*/
type DijkstraPathfinder struct {
	Maze      [][]int // Stores the entire maze for local use
	Logging   bool
	xEnd      int
	yEnd      int
	pathStack []path
}

func (this *DijkstraPathfinder) Solve() ([][]int, int) {

	borderBreaks := this.getBorderBreaks()

	if len(borderBreaks) != 2 {

		fmt.Println("The borders of the maze may only contain 2 break points")

		os.Exit(1)

	}

	this.xEnd = borderBreaks[1][0]

	this.yEnd = borderBreaks[1][1]

	this.pathStack = []path{
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

	return this.next()

}

func (this *DijkstraPathfinder) getBorderBreaks() [][]int {

	borderBreaks := [][]int{}

	for x := 0; x < len(this.Maze[0])-1; x++ {

		if this.Maze[0][x] == 0 {

			borderBreaks = append(borderBreaks, []int{x, 0})

		}

		if this.Maze[len(this.Maze)-1][x] == 0 {

			borderBreaks = append(borderBreaks, []int{x, len(this.Maze) - 1})

		}

	}

	for y := 0; y < len(this.Maze)-1; y++ {

		if this.Maze[y][0] == 0 {

			borderBreaks = append(borderBreaks, []int{0, y})

		}

		if this.Maze[y][len(this.Maze[0])-1] == 0 {

			borderBreaks = append(borderBreaks, []int{len(this.Maze[0]) - 1, y})

		}

	}

	return borderBreaks

}

func (this *DijkstraPathfinder) next() ([][]int, int) {

	if this.pathStack[0].length == -1 {

		fmt.Println("No solution could be found")

		os.Exit(1)

	}

	newPaths, endPath := this.pathStack[0].extend(this.Maze, this.xEnd, this.yEnd)

	if endPath.length == -2 {

		this.pathStack = append(this.pathStack, newPaths...)

		this.sortPathStack()

		if this.Logging {

			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()

			fmt.Println(this.pathStack)

			fmt.Print("\n\n")

		}

		return this.next()

	}

	return endPath.getPath(), len(this.pathStack)

}

// Sort the slice of paths to have the shortest path on the top - Optimal solution
func (this *DijkstraPathfinder) sortPathStack() {

	sort.Slice(this.pathStack, func(i, j int) bool {
		if this.pathStack[j].length == -1 {
			return true
		}
		if this.pathStack[i].length == -1 {
			return false
		}
		return this.pathStack[i].length < this.pathStack[j].length
	})

}

/*
	Individual paths, originating from the 'start' position
*/
type path struct {
	length int
	pathTo []node
}

// Extend the number of routes through the maze to include the nodes of the closest adjacent nodes
func (this *path) extend(maze [][]int, xEnd int, yEnd int) ([]path, path) {

	if this.length != -1 {

		adj := this.pathTo[len(this.pathTo)-1].getAdjacentNodes(maze, this.pathTo, xEnd, yEnd)

		if len(adj) > 0 {

			var newPaths []path = []path{}

			for i := 1; i < len(adj); i++ {

				newPaths = append(newPaths, path{
					length: this.length + this.pathTo[len(this.pathTo)-1].getLength(adj[i], xEnd, yEnd),
					pathTo: make([]node, len(this.pathTo)),
				})

				copy(newPaths[len(newPaths)-1].pathTo, this.pathTo)

				newPaths[len(newPaths)-1].pathTo = append(newPaths[len(newPaths)-1].pathTo, adj[i])

			}

			this.pathTo = append(this.pathTo, adj[0])

			this.length += this.pathTo[len(this.pathTo)-1].getLength(adj[0], xEnd, yEnd)

			lastNode := this.pathTo[len(this.pathTo)-1]

			endPath := path{length: -2}

			if lastNode.xPos == xEnd && lastNode.yPos == yEnd {

				endPath = *this

			}

			if endPath.length == -2 {

				for i := 0; i < len(newPaths); i++ {

					lastNode = newPaths[i].pathTo[len(this.pathTo)-1]

					if lastNode.xPos == xEnd && lastNode.yPos == yEnd {

						endPath = newPaths[i]

						break

					}

				}

			}

			return newPaths, endPath

		}

		this.length = -1

	}

	return []path{}, path{length: -2}

}

func (this *path) getPath() [][]int {

	pathTo := [][]int{}

	for _, currentNode := range this.pathTo {

		pathTo = append(pathTo, []int{currentNode.xPos, currentNode.yPos})

	}

	return pathTo

}

/*
	Junction points; one or more turns can be made
*/
type node struct {
	xPos int
	yPos int
}

// Search for nodes that can be reached on a straight path
func (this *node) getAdjacentNodes(maze [][]int, currentPath []node, xEnd int, yEnd int) []node {

	directionsToTry := this.getDirectionsToTry(len(maze)-1, currentPath)

	newNodes := []node{}

	for _, directionToMove := range directionsToTry {

		for i := 0; i < len(maze)-1; i++ {

			if this.canMove(maze, directionToMove) {

				if this.xPos+directionToMove.XDirection == xEnd && this.yPos+directionToMove.YDirection == yEnd {

					newNodes = append(newNodes, node{
						xPos: this.xPos + directionToMove.XDirection,
						yPos: this.yPos + directionToMove.YDirection,
					})

					return newNodes

				}

				if this.nodeIsPresent(maze, directionToMove) {

					newNodes = append(newNodes, node{
						xPos: this.xPos + directionToMove.XDirection,
						yPos: this.yPos + directionToMove.YDirection,
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

// Look for a possible change in direction
func (this *node) nodeIsPresent(maze [][]int, directionToMove Direction) bool {

	tempNode := node{
		xPos: this.xPos + directionToMove.XDirection,
		yPos: this.yPos + directionToMove.YDirection,
	}

	if directionToMove.XDirection != 0 {

		return tempNode.canMove(maze, Direction{
			XDirection: 0,
			YDirection: 1,
		}) || tempNode.canMove(maze, Direction{
			XDirection: 0,
			YDirection: -1,
		})

	} else {

		return tempNode.canMove(maze, Direction{
			XDirection: 1,
			YDirection: 0,
		}) || tempNode.canMove(maze, Direction{
			XDirection: -1,
			YDirection: 0,
		})

	}

}

func (this *node) canMove(maze [][]int, directionToMove Direction) bool {

	targetX := this.xPos + directionToMove.XDirection

	targetY := this.yPos + directionToMove.YDirection

	if targetY == -1 || targetX == -1 || targetY == len(maze) || targetX == len(maze[0]) {

		return false

	}

	return maze[targetY][targetX] == 0

}

func (this *node) getDirectionsToTry(mazeLength int, currentPath []node) []Direction {

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

		if this.xPos-previousNode.xPos != 0 {

			directionsToTry = append(directionsToTry, []Direction{
				Direction{0, 1},  // Down
				Direction{0, -1}, // Up
			}...)

			if this.xPos-previousNode.xPos > 0 { // If true, don't move to the left

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

			if this.yPos-previousNode.yPos > 0 { // If true, don't move upwards

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

func (this *node) getLength(adjNode node, xEnd int, yEnd int) int {
	// The absolute value of Dx + Dy
	return int(math.Abs(float64(adjNode.xPos-this.xPos))+math.Abs(float64(adjNode.xPos-this.xPos))) +
		int(math.Abs(float64(xEnd-adjNode.xPos))+math.Abs(float64(yEnd-adjNode.yPos)))

}

type Direction struct {
	XDirection int
	YDirection int
}

func (this *Direction) Increment() {

	if this.XDirection == 0 {

		if this.YDirection > 0 {

			this.YDirection += 1

		} else {

			this.YDirection -= 1

		}

	} else {

		if this.XDirection > 0 {

			this.XDirection += 1

		} else {

			this.XDirection -= 1

		}

	}

}

func (this *Direction) Decrement() bool {

	if this.XDirection == 0 {

		if this.YDirection > 0 {

			this.YDirection -= 1

		} else {

			this.YDirection += 1

		}

	} else {

		if this.XDirection > 0 {

			this.XDirection -= 1

		} else {

			this.XDirection += 1

		}

	}

	return !(this.XDirection == 0 && this.YDirection == 0)

}
