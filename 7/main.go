package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

type Sizer interface {
	Size() int
}

type File struct {
	name string
	size int
}

func (f *File) Size() int {
	return f.size
}

type Dir struct {
	name  string
	files []*File
	dirs  []*Dir
}

func (d *Dir) AddFile(f *File) {
	d.files = append(d.files, f)
}

func (d *Dir) AddDir(dir *Dir) {
	d.dirs = append(d.dirs, dir)
}

func (d *Dir) Size() int {
	size := 0
	for _, f := range d.files {
		size += f.Size()
	}
	for _, dir := range d.dirs {
		size += dir.Size()
	}
	return size
}

func newDir(name string) Dir {
	return Dir{name, make([]*File, 0), make([]*Dir, 0)}
}

type FS struct {
	files []*File
	dirs  []*Dir
}

func parseInput() []string {
	return util.NewInputFile("7").ReadLines()
}

func main() {
	flag.Parse()
	fs := FS{make([]*File, 0), make([]*Dir, 0)}
	stack := make([]*Dir, 0)

	lines := parseInput()
	i := 0
	for i < len(lines) {
		line := lines[i]
		if line[0] != '$' {
			panic("How did I get here?")
		}
		command := line[2:]
		prefix := command[0:2]
		switch prefix {
		case "cd":
			dirName := command[3:]
			if dirName == ".." {
				_, err := slices.Pop(&stack)
				util.HandleError(err)
			} else {
				if dirName == "/" {
					dir := newDir(dirName)
					fs.dirs = append(fs.dirs, &dir)
					stack = append(stack, &dir)
				} else {
					var nextDir *Dir
					curDir := stack[len(stack)-1]
					for _, d := range curDir.dirs {
						if d.name == dirName {
							nextDir = d
							break
						}
					}
					stack = append(stack, nextDir)
				}
			}
			i += 1
		case "ls":
			i += 1
			curDir := stack[len(stack)-1]
			for i < len(lines) && lines[i][0] != '$' {
				line := lines[i]
				if line[0:3] == "dir" {
					name := line[4:]
					dir := newDir(name)
					fs.dirs = append(fs.dirs, &dir)
					curDir.AddDir(&dir)
				} else {
					fileArr := strings.Split(line, " ")
					size, err := strconv.Atoi(fileArr[0])
					util.HandleError(err)
					file := File{size: size, name: fileArr[1]}
					fs.files = append(fs.files, &file)
					curDir.AddFile(&file)
				}
				i += 1
			}
		}
	}

	totalSizeOfSmallDirs := 0
	for _, dir := range fs.dirs {
		size := dir.Size()
		if size < 100000 {
			totalSizeOfSmallDirs += size
		}
	}

	freeSpace := 70000000 - fs.dirs[0].Size()
	targetFreeSpace := 30000000
	delta := targetFreeSpace - freeSpace
	currentTarget := fs.dirs[0]
	for _, dir := range fs.dirs {
		if dir.Size() > delta && dir.Size() < currentTarget.Size() {
			currentTarget = dir
		}
	}

	fmt.Println("Part 1:", totalSizeOfSmallDirs)
	fmt.Println("Part 2:", currentTarget.Size())
}
