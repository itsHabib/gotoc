package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/itsHabib/gotoc/interpreter"
)

func main() {
	f, err := getFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := gotoc(f); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func gotoc(f flags) error {
	i, err := interpreter.New(f.toDocument())
	if err != nil {
		fmt.Printf("unable to init interpreter: %s\n", err)
		os.Exit(1)
	}
	if err := i.GenerateToc(); err != nil {
		fmt.Printf("unable to generate toc: %s\n", err)
		os.Exit(1)
	}

	// just print if that's all we need to do
	if !f.write {
		fmt.Println(i.Toc())
		return nil
	}

	file, err := os.OpenFile(f.file, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("unable to open file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := write(file, i.Toc()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

// write forms the new file contents with the table of contents and writes to
// the file.
func write(file *os.File, toc string) error {
	buffer, err := addToc(file, toc)
	if err != nil {
		return fmt.Errorf("unable to generate new file contents: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("unable to seek file")
	}

	if _, err := file.Write(buffer.Bytes()); err != nil {
		return errors.New("unable to write to file")
	}

	return nil
}

func addToc(file *os.File, toc string) (*bytes.Buffer, error) {
	if file == nil {
		return nil, errors.New("can not process nil file")
	}

	// read file until we've reached the end or found the root header
	reader := bufio.NewReader(file)

	// create a buffer to hold the end file contents that contain the table of
	// contents
	buffer := new(bytes.Buffer)
	var headerFound bool
	for {
		l, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("unable to read line: %w", err)
		}

		if err == io.EOF {
			break
		}

		if !headerFound && headerDepth(strings.TrimSpace(string(l))) == 1 {
			// found root
			buffer.Write(l)
			buffer.WriteString(toc)
			headerFound = true
			continue
		}

		buffer.Write(l)
	}

	return buffer, nil
}

func headerDepth(s string) int {
	var depth int
	for ; depth < len(s) && s[depth] == '#'; depth++ {
	}

	return depth
}

func getFlags() (flags, error) {
	printFlag := flag.Bool("print", false, "only print the table of contents string and not writeFlag it to file. the table of contents are printed by default, unless -writeFlag is used.")
	writeFlag := flag.Bool("writeFlag", false, "writeFlag the table of contents string to given file. gotoc will attempt to writeFlag the contents after the root header #. If no root header is found the contents are NOT written.")
	text := flag.String("text", "", "generate a table of contents from an input string. Can not be used in conjunction with the -file flag.")
	file := flag.String("file", "", "path to a file to generate the table of contents from.")
	flag.Parse()

	var f flags
	if printFlag != nil {
		f.print = *printFlag
	}
	if writeFlag != nil {
		f.write = *writeFlag
	}
	if text != nil {
		f.text = *text
	}
	if file != nil {
		f.file = *file
	}
	if !f.write && !f.print {
		f.print = true
	}

	if err := f.validate(); err != nil {
		return flags{}, err
	}

	return f, nil
}

// flags are a wrapper for command line flags taken in
type flags struct {
	print bool
	write bool
	file  string
	text  string
}

func (f flags) validate() error {
	if f.print && f.write {
		return errors.New("can not set -write and -print flags at the same time")
	}

	if f.file != "" && f.text != "" {
		return errors.New("can not set -file and -text flags t the same time")
	}

	if f.file == "" && f.text == "" {
		return errors.New("must set one of -file or -text flags")
	}

	return nil
}

func (f flags) toDocument() interpreter.Document {
	return interpreter.Document{
		Content: f.text,
		Path:    f.file,
	}
}
