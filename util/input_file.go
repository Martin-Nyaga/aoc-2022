package util

import (
	"bufio"
	"flag"
	"os"
)

var UseSampleInput = flag.Bool("sample", false, "Use sample input")

type InputFile struct {
	SampleFilePath string
	InputFilePath  string
}

func NewInputFile(problem string) InputFile {
	return InputFile{
		SampleFilePath: "./" + problem + "/sample.txt",
		InputFilePath:  "./" + problem + "/input.txt",
	}
}

func (i InputFile) filePath() string {
	if *UseSampleInput {
		return i.SampleFilePath
	} else {
		return i.InputFilePath
	}
}

func (i InputFile) ReadLines() []string {
	f, err := os.Open(i.filePath())
	HandleError(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	result := make([]string, 0)
	for s.Scan() {
		result = append(result, s.Text())
	}

	return result
}

func (i InputFile) ReadBytes() []byte {
	s, err := os.ReadFile(i.filePath())
	HandleError(err)
	return s
}

func (i InputFile) ReadToString() string {
	return string(i.ReadBytes())
}
