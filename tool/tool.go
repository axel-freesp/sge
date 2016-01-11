package tool

import (
	"fmt"
	"os"
	"strings"
)

var Debug bool = false
var VerboseErr bool = true

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AbsInt(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func ReadFile(filepath string) (data []byte, err error) {
	var file *os.File
	var size int64
	var count int

	if Debug {
		fmt.Println("Reading from", filepath)
	}
	file, err = os.Open(filepath)
	if err != nil {
		if VerboseErr {
			fmt.Println("Failed to open file", filepath)
		}
		return
	}
	defer file.Close()

	size, err = file.Seek(0, os.SEEK_END)
	if err != nil {
		if VerboseErr {
			fmt.Println("Failed to seek file", filepath)
		}
		return
	}
	file.Seek(0, os.SEEK_SET)

	data = make([]byte, size)
	count, err = file.Read(data)
	if err != nil {
		if VerboseErr {
			fmt.Println("Failed to read file", filepath)
		}
		return
	}

	if count != int(size) {
		if VerboseErr {
			fmt.Println("Incomplete read of file", filepath)
		}
		return
	}
	return
}

func WriteFile(filepath string, data []byte) (err error) {
	var file *os.File
	if Debug {
		fmt.Println("Writing to ", filepath)
	}
	file, err = os.Create(filepath)
	if err != nil {
		if VerboseErr {
			fmt.Println("Failed to create file ", filepath)
		}
		return
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		if VerboseErr {
			fmt.Println("Failed to write ", len(data), " bytes to ", filepath)
		}
		return
	}
	return
}

func Basename(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return path
	}
	return path[idx+1:]
}

func Dirname(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return "."
	}
	return path[:idx]
}

func IsSubPath(short, long string) bool {
	return short == long[:len(short)]
}

func RelPath(short, long string) string {
	if !IsSubPath(short, long) {
		return long
	}
	return long[len(short)+1:]
}

func Suffix(path string) string {
	idx := strings.LastIndex(path, ".")
	if idx < 0 {
		return path
	}
	return path[idx+1:]
}
