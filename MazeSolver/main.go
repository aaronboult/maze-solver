package main

import (
	"dijkstrapathfinder"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {

	var outputName, mazeName string

	var showNumericalSolution bool

	var preventOutputMaze bool

	flag.StringVar(&mazeName, "maze", "maze.png", "The file name for the maze to solve. Defaults to 'maze.png'")

	flag.StringVar(&outputName, "output", "maze_output.png", "The file name for the output solution. Defaults to 'maze_output.png'")

	flag.BoolVar(&showNumericalSolution, "showNumericalSolution", false, "Whether or not to show the integer representation of the solution (May be a large output for larger mazes)")

	flag.BoolVar(&preventOutputMaze, "preventWritingOutput", false, "Prevents writing an output maze with the solution")

	flag.Parse()

	fmt.Printf("Reading maze: %s\n", mazeName)

	maze := loadMaze(mazeName)

	solver := dijkstrapathfinder.DijkstraPathfinder{Maze: maze, Logging: showNumericalSolution} // Logging dictates whether to show the live path stack while pathfinding

	fmt.Println("Solving maze using Dijkstra's Pathfinder algorithm...")

	solution, numberOfPaths := solver.Solve()

	fmt.Printf("Tried %d paths.\n", numberOfPaths)

	if showNumericalSolution {

		fmt.Printf("Solution: %v\n", solution)

	}

	if !preventOutputMaze {

		solutionImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{len(maze[0]), len(maze)}}) // Create a new blank image the same size as the given maze

		for y := 0; y < len(maze); y++ {

			for x := 0; x < len(maze[0]); x++ {

				if maze[y][x] == 1 {

					solutionImage.Set(x, y, color.Black)

				} else {

					solutionImage.Set(x, y, color.White)

				}

			}

		} // Copy the given maze image to the newly created image

		red := color.RGBA{0xff, 0, 0, 0xff}

		for index, current := range solution {

			solutionImage.Set(current[0], current[1], red)

			if index != len(solution)-1 {

				direction := dijkstrapathfinder.Direction{
					XDirection: solution[index+1][0] - current[0],
					YDirection: solution[index+1][1] - current[1],
				}

				for direction.Decrement() {

					solutionImage.Set(current[0]+direction.XDirection, current[1]+direction.YDirection, red)

				}

			}

		} // Draw the red path line on the newly created image

		fmt.Printf("Creating maze solution: %s\n", outputName)

		newFile, _ := os.Create(outputName)

		png.Encode(newFile, solutionImage)

		newFile.Close()

	}

}

func loadMaze(mazeName string) [][]int {

	img := loadImage("./" + mazeName)

	var maze [][]int = [][]int{}

	var r uint32 // The red value in the image - used to determine whether a pixel is black or white
				 // 0 represents white, 1 represents black

	for rowIndex := 0; rowIndex < img.Bounds().Dy(); rowIndex++ {

		maze = append(maze, []int{})

		for columnIndex := 0; columnIndex < img.Bounds().Dx(); columnIndex++ {

			r, _, _, _ = img.At(columnIndex, rowIndex).RGBA()

			if r == 0 {

				maze[rowIndex] = append(maze[rowIndex], 1)

				continue

			}

			maze[rowIndex] = append(maze[rowIndex], 0)

		}

	}

	return maze

}

func loadImage(filePath string) image.Image {

	file, err := os.Open(filePath)

	defer file.Close()

	checkError(err)

	img, _, err := image.Decode(file)

	checkError(err)

	return img

}

func checkError(err error) {

	if err != nil {

		panic(err)

	}

}
