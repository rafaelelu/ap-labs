package main

import (
	"fmt"
	"flag"
	"log"
	"time"
	"math/rand"
	"strconv"
	"bufio"
	"os"
	"github.com/golang-collections/collections/stack" // go get github.com/golang-collections/collections/stack
)

// type of cells
const STREET = 0
const BUILDING = 1

// directions
const NO_DIR = -1
const LEFT = 0
const RIGHT = 1
const DOWN = 2
const UP = 3
const LDOWN = 4
const RDOWN = 5
const LUP = 6
const RUP = 7

type cell struct {
	x int
	y int
	typeOfCell int
	dir int
	hasCar bool
	greenLight bool
}

type car struct {
	id int
	x int
	y int
	speed int
	path []cell
	idle int
}

// this struct is used in the BFS when searching for a path
type qElement struct {
	previous *qElement
	c cell
}

type semaphore struct {
	cells []cell
	index int
	speed int
}

var width int
var nCars int
var nSemaphores int
var r *rand.Rand
var board [][]cell
var paths [][]cell
var cars []car

func main() {

	widthFlag := flag.Int("width", 9, "width needs to be a positive integer")
	nCarsFlag := flag.Int("cars", 9, "cars needs to be a positive integer")
	nSemFlag  := flag.Int("semaphores", 4, "semaphores needs to be a positive integer")
	flag.Parse();
	width = *widthFlag
	nCars = *nCarsFlag
	nSemaphores = *nSemFlag

	validateWidth()
	
	r = rand.New(rand.NewSource(time.Now().UnixNano())) // seed
	var streetCells []cell
	board, streetCells = createBoard()

	validateCars(width)
	
	intersections := getIntersections()

	validateSemaphores(len(intersections))

	createSemaphores(intersections)

	ch := make(chan int, nCars)
	createCars(streetCells, &ch)
	showMap(&ch)

	fmt.Print("Do you want to see each car's route? (y/n): ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

	for ;input.Text() != "y" && input.Text() != "n"; {
		fmt.Print("Do you want to see each car's route? (y/n): ")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
	}

	if input.Text() == "y" {
		printPaths()
	}

}

func validateWidth() {
	if width < 9 || (width-2)%7 != 0 {
		log.Fatalf("ERROR: The width must be (a multiple of 7) + 2. For example: 9, 16, 30, 44, 51, ...")
	}
}

func validateCars(width int) {
	if nCars > width || nCars < 0 {
		log.Fatalf("ERROR: With a width of %v the number of cars must be a positive integer less or equal than %v",width, width)
	}
}

func validateSemaphores(nIntersections int) {
	if nSemaphores < 0 || nSemaphores > nIntersections {
		log.Fatalf("ERROR: With a width of %v, you have %v intersections, so semaphores must be a positive integer less or equal than %v", width, nIntersections, nIntersections)
	}
}

func createBoard() ([][]cell, []cell) {
	
	var board = make([][]cell, 0)
	var streetCells = make([]cell, 0)

	for i := 0; i < width; i++ {
		var line = make([]cell, 0)
		iMod := i%7
		for j := 0; j < width; j++ {
			jMod := j%7
			if iMod < 2 || jMod < 2 {
				dir := NO_DIR
				if iMod == 0 && jMod == 0 {
					dir = LDOWN
				} else if iMod == 1 && jMod == 0 {
					dir = RDOWN
				} else if iMod == 0 && jMod == 1 {
					dir = LUP
				} else if iMod == 1 && jMod == 1 {
					dir = RUP
				} else if iMod == 0 {
					dir = LEFT
				} else if iMod == 1 {
					dir = RIGHT
				} else if jMod == 0 {
					dir = DOWN
				} else if jMod == 1 {
					dir = UP
				}
				c := cell{i, j, STREET, dir, false, true}
				if dir == LEFT || dir == RIGHT || dir == UP || dir == DOWN {
					streetCells = append(streetCells, c)
				}
				line = append(line, c) // That cell is part of a street
			} else {
				c := cell{i, j, BUILDING, NO_DIR, false, true}
				line = append(line, c) // That cell is part of a building
			}
		}
		board = append(board, line)
	}
	return board, streetCells

}

func getIntersections() [][]cell {
	intersections := make([][]cell, 0)
	for i := 0; i < width; i += 7 {
		for j := 0; j < width; j += 7 {
			intersection := make([]cell, 0)
			if i > 0 {
				intersection = append(intersection, board[i-1][j])
			}
			if j > 0 {
				intersection = append(intersection, board[i+1][j-1])
			}
			if i < width - 2 {
				intersection = append(intersection, board[i+2][j+1])
			}
			if j < width - 2 {
				intersection = append(intersection, board[i][j+2])
			}
			intersections = append(intersections, intersection)
		}
	}
	return intersections
}

func createSemaphores(intersections [][]cell) {
	semaphores := make([]semaphore, 0)

	for i := 0; i < nSemaphores; i++ {
		n := r.Intn(len(intersections))

		intersection := intersections[n]

		intersections[len(intersections)-1], intersections[n] = intersections[n], intersections[len(intersections)-1]
		intersections = intersections[:len(intersections)-1]

		speed := r.Intn(1200 - 800) + 800
		s := semaphore{intersection, 0, speed}
		setInitState(&s)
		semaphores = append(semaphores, s)
	}

	for i := 0; i < len(semaphores); i++ {
		index := i
		go func() {
			for ;; {
				changeState(&semaphores[index])
				time.Sleep(time.Duration(semaphores[index].speed) * time.Millisecond)
			}
		} ()
	}


}

func setInitState(s *semaphore) {

	for i := 0; i < len(s.cells); i++ {
		x := s.cells[i].x
		y := s.cells[i].y
		board[x][y].greenLight = false
	}
}

func changeState(s *semaphore) {
	length := len(s.cells)
	cX := s.cells[s.index].x
	cY := s.cells[s.index].y
	board[cX][cY].greenLight = false
	s.index = (s.index+1)%length
	nX := s.cells[s.index].x
	nY := s.cells[s.index].y
	board[nX][nY].greenLight = true
}

func printBoard() {

	fmt.Println("\033[H\033[2J")

	for i := 0; i < width; i++ {
		line := ""
		for j := 0; j < width; j++ {
			if board[i][j].hasCar {
				line += " ■ "
			} else if !board[i][j].greenLight {
				d := board[i][j].dir
				if d == LEFT || d == RIGHT {
					line += "---"
				}
				if d == DOWN || d == UP {
					line += " | "
				}
			} else {
				switch(board[i][j].typeOfCell) {
				case STREET:
					switch(board[i][j].dir) {
					case LEFT:
						line += "---"
					case RIGHT:
						line += "---"
					case DOWN:
						line += " | "
					case UP:
						line += " | "
					case LDOWN:
						line += "  \\"
					case RDOWN:
						line += "  /"
					case LUP:
						line += "/  "
					case RUP:
						line += "\\  "
					default:
						line += "   "
					}
					break;
				case BUILDING:
					line += "▒▒▒"
				}
			}
		}
		if i < nCars {
			cIndex := strconv.Itoa(i)
			if i < 10 {
				cIndex = "0" + cIndex
			}

			if cars[i].speed == 0 {
				line += "	Car #" + cIndex + ": 		FINISHED!"
			} else {
				if cars[i].speed > 240 {
					line += "	Car #" + cIndex + "'s speed: 0 km/h"
				} else {
					line += "	Car #" + cIndex + "'s speed: " + strconv.Itoa(1250 / cars[i].speed) + "km/h"
				}
			}
		}
		fmt.Println(line);
	}

}

func createCars(streetCells []cell, ch *chan int) {
	cars = make([]car, 0)
	for i := 0; i < nCars; i++ {
		showLoadingScreen(i)
		n := r.Intn(len(streetCells))
		cell1 := streetCells[n]

		streetCells[len(streetCells)-1], streetCells[n] = streetCells[n], streetCells[len(streetCells)-1]
		streetCells = streetCells[:len(streetCells)-1]

		n2 := r.Intn(len(streetCells))
		cell2 := streetCells[n2]
		
		speed := r.Intn(250 - 50) + 50
		
		path := getPath(cell1, cell2)
		paths = append(paths, path)

		c := car{i, cell1.x, cell1.y, speed, path, 0}


		cars = append(cars, c)
		addCar(c)
	}

	for i := 0; i < len(cars); i++ {
		index := i
		go func() {
			for ; len(cars[index].path) > 0; {
				time.Sleep(time.Duration(cars[index].speed) * time.Millisecond)
				moveCar(&cars[index])
			}
			*ch <- cars[index].id
			cars[index].speed = 0
			removeCar(&cars[index])
		} ()
	}
}

func showLoadingScreen(createdCars int) {
	fmt.Println("\033[H\033[2J")
	fmt.Println()
	fmt.Println("                                                                                                                                                         ")
	fmt.Println("        ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄         ▄       ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄        ")
	fmt.Println("       ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░▌       ▐░▌     ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▌       ")
	fmt.Println("       ▐░█▀▀▀▀▀▀▀▀▀ ▀▀▀▀█░█▀▀▀▀ ▀▀▀▀█░█▀▀▀▀▐░▌       ▐░▌      ▀▀▀▀█░█▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀▀▀▐░█▀▀▀▀▀▀▀▀▀ ▀▀▀▀█░█▀▀▀▀▐░█▀▀▀▀▀▀▀▀▀        ")
	fmt.Println("       ▐░▌              ▐░▌         ▐░▌    ▐░▌       ▐░▌          ▐░▌    ▐░▌       ▐░▐░▌       ▐░▐░▌         ▐░▌              ▐░▌    ▐░▌                 ")
	fmt.Println("       ▐░▌              ▐░▌         ▐░▌    ▐░█▄▄▄▄▄▄▄█░▌          ▐░▌    ▐░█▄▄▄▄▄▄▄█░▐░█▄▄▄▄▄▄▄█░▐░█▄▄▄▄▄▄▄▄▄▐░█▄▄▄▄▄▄▄▄▄     ▐░▌    ▐░▌                 ")
	fmt.Println("       ▐░▌              ▐░▌         ▐░▌    ▐░░░░░░░░░░░▌          ▐░▌    ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▌    ▐░▌    ▐░▌                 ")
	fmt.Println("       ▐░▌              ▐░▌         ▐░▌     ▀▀▀▀█░█▀▀▀▀           ▐░▌    ▐░█▀▀▀▀█░█▀▀▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀▀▀▐░█▀▀▀▀▀▀▀▀▀     ▐░▌    ▐░▌                 ")
	fmt.Println("       ▐░▌              ▐░▌         ▐░▌         ▐░▌               ▐░▌    ▐░▌     ▐░▌ ▐░▌       ▐░▐░▌         ▐░▌              ▐░▌    ▐░▌                 ")
	fmt.Println("       ▐░█▄▄▄▄▄▄▄▄▄ ▄▄▄▄█░█▄▄▄▄     ▐░▌         ▐░▌               ▐░▌    ▐░▌      ▐░▌▐░▌       ▐░▐░▌         ▐░▌          ▄▄▄▄█░█▄▄▄▄▐░█▄▄▄▄▄▄▄▄▄        ")
	fmt.Println("       ▐░░░░░░░░░░░▐░░░░░░░░░░░▌    ▐░▌         ▐░▌               ▐░▌    ▐░▌       ▐░▐░▌       ▐░▐░▌         ▐░▌         ▐░░░░░░░░░░░▐░░░░░░░░░░░▌       ")
	fmt.Println("        ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀      ▀           ▀                 ▀      ▀         ▀ ▀         ▀ ▀           ▀           ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀        ")
	fmt.Println("                                                                                                                                                         ")
	fmt.Println("                    ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄       ▄▄ ▄         ▄ ▄           ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄ ▄▄▄▄▄▄▄▄▄▄▄                          ")
	fmt.Println("                   ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░▌     ▐░░▐░▌       ▐░▐░▌         ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░░░░░░░░░░░▌                         ")
	fmt.Println("                   ▐░█▀▀▀▀▀▀▀▀▀ ▀▀▀▀█░█▀▀▀▀▐░▌░▌   ▐░▐░▐░▌       ▐░▐░▌         ▐░█▀▀▀▀▀▀▀█░▌▀▀▀▀█░█▀▀▀▀▐░█▀▀▀▀▀▀▀█░▐░█▀▀▀▀▀▀▀█░▌                         ")
	fmt.Println("                   ▐░▌              ▐░▌    ▐░▌▐░▌ ▐░▌▐░▐░▌       ▐░▐░▌         ▐░▌       ▐░▌    ▐░▌    ▐░▌       ▐░▐░▌       ▐░▌                         ")
	fmt.Println("                   ▐░█▄▄▄▄▄▄▄▄▄     ▐░▌    ▐░▌ ▐░▐░▌ ▐░▐░▌       ▐░▐░▌         ▐░█▄▄▄▄▄▄▄█░▌    ▐░▌    ▐░▌       ▐░▐░█▄▄▄▄▄▄▄█░▌                         ")
	fmt.Println("                   ▐░░░░░░░░░░░▌    ▐░▌    ▐░▌  ▐░▌  ▐░▐░▌       ▐░▐░▌         ▐░░░░░░░░░░░▌    ▐░▌    ▐░▌       ▐░▐░░░░░░░░░░░▌                         ")
	fmt.Println("                    ▀▀▀▀▀▀▀▀▀█░▌    ▐░▌    ▐░▌   ▀   ▐░▐░▌       ▐░▐░▌         ▐░█▀▀▀▀▀▀▀█░▌    ▐░▌    ▐░▌       ▐░▐░█▀▀▀▀█░█▀▀                          ")
	fmt.Println("                             ▐░▌    ▐░▌    ▐░▌       ▐░▐░▌       ▐░▐░▌         ▐░▌       ▐░▌    ▐░▌    ▐░▌       ▐░▐░▌     ▐░▌                           ")
	fmt.Println("                    ▄▄▄▄▄▄▄▄▄█░▌▄▄▄▄█░█▄▄▄▄▐░▌       ▐░▐░█▄▄▄▄▄▄▄█░▐░█▄▄▄▄▄▄▄▄▄▐░▌       ▐░▌    ▐░▌    ▐░█▄▄▄▄▄▄▄█░▐░▌      ▐░▌                          ")
	fmt.Println("                   ▐░░░░░░░░░░░▐░░░░░░░░░░░▐░▌       ▐░▐░░░░░░░░░░░▐░░░░░░░░░░░▐░▌       ▐░▌    ▐░▌    ▐░░░░░░░░░░░▐░▌       ▐░▌                         ")
	fmt.Println("                    ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀         ▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀▀▀▀ ▀         ▀      ▀      ▀▀▀▀▀▀▀▀▀▀▀ ▀         ▀                          ")
	fmt.Println("                                                                                                                                                         ")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("                                                                                                            ")
	/*fmt.Println("                                               ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄  ")


	barLine := "                                              ▐"
	barProgress := loadingProgress(59, nCars, createdCars)
	for i := 0; i < 59; i++ {
		if i < barProgress {
			barLine += "░"
		} else {
			barLine += " "
		}
	}
	barLine += "▌"
	fmt.Println(barLine)
	fmt.Println("                                               ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀  ")*/
	fmt.Println()
	fmt.Println()
	fmt.Println()
	traveledDistance := progressString(loadingProgress(80, nCars, createdCars))
	fmt.Println( traveledDistance + "                      ____________________                                 ")
	fmt.Println( traveledDistance + "                    //|           |        \\                              ")
	fmt.Println( traveledDistance + "                  //  |           |          \\                            ")
	fmt.Println( traveledDistance + "     ___________//____|___________|__________()\\__________________        ")
	fmt.Println( traveledDistance + "   /__________________|_=_________|_=___________|_________________{}       ")
	fmt.Println( traveledDistance + "   [           ______ |           | .           | ==  ______      { }      ")
	fmt.Println( traveledDistance + " __[__        /##  ##\\|           |             |    /##  ##\\    _{# }_  ")
	fmt.Println( traveledDistance + "{_____)______|##    ##|___________|_____________|___|##    ##|__(______}   ")
	fmt.Println( traveledDistance + "            /  ##__##                              /  ##__##        \\     ")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	// 80
}

func progressString(nProgress int) string {
	s := ""
	for i := 0; i < nProgress; i++ {
		s += " "
	}
	return s
}

func loadingProgress(barSize int, totalCars int, cars int) int {
	f := (barSize * cars) / totalCars
	return f
}

func addCar(c car) {
	i := c.x
	j := c.y
	if !board[i][j].hasCar {
		board[i][j].hasCar = true
	}
}

func showMap(ch *chan int) {
	for ;; {
		printBoard()
		fmt.Printf("Progress: %v / %v\n",len(*ch), nCars)
		if len(*ch) >= nCars {
			break
		}
		time.Sleep(41 * time.Millisecond)
	}
}

func removeCar(c *car) {
	i := c.x
	j := c.y
	if board[i][j].hasCar {
		board[i][j].hasCar = false
	}
}

func moveCar(c *car) {
	cX := c.x
	cY := c.y
	nextCell := c.path[0]
	nX := nextCell.x
	nY := nextCell.y
	if !board[nX][nY].hasCar && board[cX][cY].greenLight {
		board[cX][cY].hasCar = false
		c.x = nX
		c.y = nY
		board[nX][nY].hasCar = true
		c.path = c.path[1:]
		if c.speed > 50 {
			c.speed -= 10
		}
		c.idle = 0
	} else {
		if c.idle <= 2 {
			c.idle++
			if c.speed < 250 {
				c.speed += 10
			}
		} else {
			c.speed = 250
		}
	}
}

func getPath(source cell, destination cell) []cell {
	q := qElement{nil, source}
	visited := make([]cell, 0)
	queue := make([]qElement, 0)
	queue = append(queue, q)
	for ; len(queue) != 0; {
		curr := queue[0]
		queue = queue[1:]
		if curr.c == destination {
			return buildpath(&curr)
		}
		visited = append(visited, curr.c)
		neighbors := getNeighbors(visited, curr.c)
		for i := 0; i < len(neighbors); i++ {
			q2 := qElement{&curr, neighbors[i]}
			if !beenVisited(visited, q2.c) {
				queue = append(queue, q2)
			}
		}
	}
	return nil
}

func getNeighbors(visited []cell, source cell) []cell {
	var neighbors = make([]cell, 0)
	x := source.x
	y := source.y

	d := source.dir
	if d == LEFT || d == LDOWN || d == LUP {
		if y > 0 {
			c := board[x][y-1]
			if !beenVisited(visited, c) && c.typeOfCell != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	}

	if d == DOWN || d == LDOWN || d == RDOWN {
		if x < width - 1 {
			c := board[x+1][y]
			if !beenVisited(visited, c) && c.typeOfCell != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	}

	if d == RIGHT || d == RDOWN || d == RUP {
		if y < width - 1 {
			c := board[x][y+1]
			if !beenVisited(visited, c) && c.typeOfCell != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	}

	if d == UP || d == LUP || d == RUP {
		if x > 0 {
			c := board[x-1][y]
			if !beenVisited(visited, c) && c.typeOfCell != BUILDING {
				neighbors = append(neighbors, c)
			}
		}
	}

	return neighbors
}

func beenVisited(visited []cell, c cell) bool {
	length := len(visited)
	for i := 0; i < length; i++ {
		if visited[i] == c {
			return true
		}
	}
	return false
}

func buildpath(q *qElement) []cell {
	path := make([]cell,0)
	s := stack.New()
	curr := q
	for ; curr != nil; {
		s.Push(curr.c)
		curr = curr.previous
	}
	s.Pop()
	for ; s.Len() > 0; {
		path = append(path, s.Pop().(cell))
	}
	return path
}

func printPaths() {
	for i := 0; i < len(paths); i++ {
		fmt.Printf("Car #%d's route: \n", i)
		for j := 0; j < len(paths[i]); j++ {
			fmt.Printf(" -> (%d,%d)", paths[i][j].x, paths[i][j].y)
		}
		fmt.Println()
		fmt.Println()
	}
}

func printCell(c cell) {
	fmt.Println("x:", c.x, ", y:", c.y)
}

func printCellSlice(slice []cell) {
	length := len(slice)
	for i := 0; i < length; i++ {
		printCell(slice[i])
	}
}