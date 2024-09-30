package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

type tree struct {
	level    int
	isDir    bool
	path     string
	size     int64
	fullPath string
}

func countCharacters(s string) int {
	count := 0

	for _, char := range s {
		if char == '\\' {
			count++
		}
	}

	return count
}

func dirTree(output io.Writer, Hpath string, printFiles bool) error {
	// Вызовите функцию для вывода содержимого директории
	slice := make([]tree, 0)
	countMap := make(map[string]int)
	levels := make(map[int]bool, 10)

	err := filepath.WalkDir(Hpath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == Hpath {
			return nil
		}

		if info.IsDir() {
			if _, exists := countMap[filepath.Dir(path)]; exists {
				countMap[filepath.Dir(path)]++
			} else {
				countMap[filepath.Dir(path)] = 1
			}
			treeExmple := tree{countCharacters(path) - 1, true, filepath.Base(path), 0, path}
			slice = append(slice, treeExmple)
		} else if printFiles && !info.IsDir() {
			if _, exists := countMap[filepath.Dir(path)]; exists {
				countMap[filepath.Dir(path)]++
			} else {
				countMap[filepath.Dir(path)] = 1
			}
			fileInfo, err := info.Info()
			if err != nil {
				return err
			}
			treeExmple := tree{countCharacters(path) - 1, false, filepath.Base(path), fileInfo.Size(), path}
			slice = append(slice, treeExmple)
		}
		return nil
	})
	var formattedString string

	for _, value := range slice {
		level := ""
		for i := 0; i < value.level; i++ {
			if levels[i] {
				level += "	"
			} else {
				level += "│	"
			}
		}
		var baseSign string

		if val := countMap[filepath.Dir(value.fullPath)]; val > 1 {
			countMap[filepath.Dir(value.fullPath)]--
			baseSign = "├───"
			levels[value.level] = false
		} else {
			baseSign = "└───"
			levels[value.level] = true
		}
		if value.isDir {
			formattedString += level + baseSign + value.path + "\n"
		} else {
			var size string
			if value.size == 0 {
				size = " (empty)"
			} else {
				size = fmt.Sprintf(" (%db)", value.size)
			}
			formattedString += level + baseSign + value.path + size + "\n"
		}
	}

	// fmt.Println(formattedString)
	fmt.Fprintf(output, formattedString)
	if err != nil {
		fmt.Printf("Ошибка при обходе директории: %v\n", err)
	}
	return nil
}
