package main

import ("math"
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
	)

func main() {
	n, error := strconv.Atoi(os.Args[1])
	fmt.Println("- Generating a [", n, "] sides figure")
	fmt.Println("- Figure's vertices")
	fmt.Println(GetPerimeter(GetFigure(n)))
	_ = error
}

func GetFigure(n int) Lines {

	var path Path
	for i := 0; i < n; i++ {
		p := Point{RandInRange(-100, 100), RandInRange(-100, 100)}
		path = append(path, p)
	}
	var lines Lines
	lines = GetLines(path)

	var lc LineComb
	lc = GetComb(lines)

	for i := 0; i < len(lc); i++ {
		if(IsValidFigure(lc[i])) {
			PrintVertices(lc[i])
			return lc[i]
		}
	}
	PrintVertices(lines)
	return lines
}

func PrintVertices(lines Lines) {
	for i := 0; i < len(lines); i++ {
		fmt.Println("- (", lines[i].p1.X(), ",",lines[i].p1.Y(),")")
	}
}

func GetPerimeter(lines Lines) float64 {
	p := 0.0
	fmt.Println("- Figure's Perimeter")
	fmt.Print("- ")
	for i := 0; i < len(lines); i++ {
		d := Distance(lines[i].p1, lines[i].p2)
		fmt.Print(d)
		if (i < len(lines)-1) {
			fmt.Print(" + ")
		}
		p = p + d
	}
	fmt.Print(" = ")
	return p
}

func IsValidFigure(lines Lines) bool {
	for i := 0; i < len(lines); i++ {
		for j := i; j < len(lines); j++ {
			if(DoIntersect(lines[i].p1, lines[i].p2, lines[j].p1, lines[j].p2)) {
				return false;
			}
		}
	}
	return true
}

func GetLines(path Path) Lines {
	var lines Lines
	for i := 0; i < len(path); i++ {
		l := Line{path[i], path[(i+1) % len(path)]}
		lines = append(lines, l)
	}
	return lines
}

func GetComb(lines Lines) LineComb {
	var lc LineComb
	for i := 0; i < len(lines); i++ {
		Swap(&lines[i], &lines[(i+1) % len(lines)])
		lc = append(lc, lines)
		Swap(&lines[i], &lines[(i+1) % len(lines)])
	}
	return lc
}

func Swap(l1, l2 *Line) {
	tmp := l1
	l1 = l2
	l2 = tmp
}

func OnSegment(p, q, r Point) bool {
	return (q.X() <= Max(p.X(), r.X()) && q.X() >= Min(p.X(), r.X()) && q.Y() <= Max(p.Y(), r.Y()) && q.Y() >= Max(p.Y(), r.Y()))
}


func Max(a, b float64)  float64 {
	if (a > b) {
		return a;
	}
	return b;
}

func Min(a, b float64) float64 {
	if (a < b) {
		return a;
	}
	return b;
}

func Orientation(p, q, r Point) int {
	val := (q.Y() - p.Y() * (r.X() - q.X()) - (q.X() - p.X()) * (r.Y() - q.Y()))
	min := -0.0001
	max := 0.0001
	if (val > min && val < max) {
		return 0
	}
	if (val > max) {
		return 1
	} else {
		return 2
	}
}

func DoIntersect(p1, q1, p2, q2 Point) bool {
	o1 := Orientation(p1, q1, p2)
	o2 := Orientation(p1, q1, q2)
	o3 := Orientation(p2, q2, p1)
	o4 := Orientation(p2, q2, q1)

	if (o1 != o2 && o3 != o4) {
		return true
	}

	if (o1 == 0 && OnSegment(p1, p2, q1)) {
		return true
	}
	if (o2 == 0 && OnSegment(p1, q2, q1)) {
		return true
	}
	if (o3 == 0 && OnSegment(p2, p1, q2)) {
		return true
	}
	if (o4 == 0 && OnSegment(p2, q1, q2)) {
		return true
	}
	return false
}

func RandInRange(min int, max int) float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	return r.Float64()* float64((max - min) + min)
}


type Point struct{ x, y float64 }
type Line struct{p1, p2 Point}
type Lines []Line
type Path []Point
type LineComb []Lines


// traditional function
func Distance(p, q Point) float64 {
        return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
        return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

func (p Point) X() float64 {
	return p.x
}

func (p Point) Y() float64 {
	return p.y
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.


// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
        sum := 0.0
        for i := range path {
                if i > 0 {
                        sum += path[i-1].Distance(path[i])
                }
        }
        return sum
}

//!-path
