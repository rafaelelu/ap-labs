package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

        if len(os.Args) < 2 {
                fmt.Println("Usage: ./dir-scan <directory>")
                os.Exit(1)
        }
        scanDir(os.Args[1])
}

func scanDir(dir string) error {
	directories := 0
	symLinks := 0
	devices := 0
	sockets := 0
	other := 0
	err := filepath.Walk(dir, func(dir string, fileI os.FileInfo, err error) error{
		fileI, err =  os.Lstat(dir)
		if (fileI.Mode() & os.ModeDir != 0){
			directories = directories + 1
		} else if (fileI.Mode() & os.ModeSymlink != 0){
			symLinks = symLinks + 1
		} else if(fileI.Mode() & os.ModeDevice != 0){
			devices = devices + 1
		} else if(fileI.Mode() & os.ModeSocket != 0){
			sockets = sockets + 1
		} else {
			other = other + 1
		}
		return err
	})
	fmt.Println("Directory Scanner Tool")
	fmt.Println("+-------------------------+------+")
	fmt.Println("| Path                    |",dir,"|")
	fmt.Println("+-------------------------+------+")
	fmt.Println("| Directories             | ",directories," |")
	fmt.Println("| Symbolic Links          | ",symLinks," |")
	fmt.Println("| Devices                 | ",devices," |")
	fmt.Println("| Sockets                 | ",sockets," |")
	fmt.Println("| Other files             | ",other," |")
	fmt.Println("+-------------------------+------+")
	return err
}
