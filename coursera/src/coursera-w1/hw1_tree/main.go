package main

import (
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type sortable []os.FileInfo

func (a sortable) Len() int {
	return len(a)
}

func (a sortable) Swap(i, j int) {
	a[j], a[i] = a[i], a[j]
}

func (a sortable) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	processDir(out, path, printFiles, 0, 0)
	return nil
}

func processDir(out io.Writer, path string, printFiles bool, depth, lastBefore int) {
	counter := 0
	isLast := false

	files := getFilesToProcess(path)
	filesCount := getFilesCount(files, printFiles)

	for _, file := range files {
		if printFiles || file.IsDir() {
			counter++
			isLast = counter == filesCount
			out.Write([]byte(makeLine(depth, isLast, lastBefore) + makeFileInfo(file, printFiles) + "\n"))
		}

		if file.IsDir() {
			if isLast {
				lastBefore++
			}
			processDir(out, makeNewPath(path, file.Name()), printFiles, depth+1, lastBefore)
		}
	}
}

func makeFileInfo(file os.FileInfo, printFiles bool) string {
	result := file.Name()
	size := "empty"

	if printFiles && !file.IsDir() {
		if file.Size() > 0 {
			size = strconv.FormatInt(file.Size(), 10) + "b"
		}

		result += " " + "(" + size + ")"
	}

	return result
}

func makeLine(depth int, isLast bool, lastBefore int) string {
	treeSymbol := "├"

	if isLast {
		treeSymbol = "└"
	}

	r := ""
	r = "───" + r
	r = treeSymbol + r

	for i := 0; i < depth; i++ {
		if i < lastBefore {
			r = "\t" + r
		} else {
			r = "│\t" + r
		}
	}

	return r
}

func makeNewPath(currentPath, addToPath string) string {
	return strings.Join([]string{currentPath, addToPath}, string(os.PathSeparator))
}

func getFilesToProcess(path string) sortable {
	dir, _ := os.Open(path)
	dirFiles, _ := dir.Readdir(0)
	sortableDirFiles := sortable(dirFiles)
	sort.Sort(sortableDirFiles)

	return sortableDirFiles
}

func getFilesCount(files sortable, printFiles bool) int {
	if printFiles {
		return len(files)
	}

	count := 0
	for _, file := range files {
		if file.IsDir() {
			count++
		}
	}

	return count
}
