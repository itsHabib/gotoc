package interpreter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/itsHabib/gotoc/parser"
)

const (
	generatedComment = "<!-- gotoc generated table of contents -->"
)

type Interpreter struct {
	document    Document
	headerCount map[string]int
	indent      int
	parser      *parser.Parser
	toc         string
}

func NewInterpreter(document Document) (*Interpreter, error) {
	if document.Path != "" {
		if _, err := os.Stat(document.Path); os.IsNotExist(err) {
			return nil, errors.New("no file found at document path")
		}
	}

	i := Interpreter{
		parser:      parser.NewParser(),
		headerCount: make(map[string]int),
		document:    document,
	}

	if err := i.validate(); err != nil {
		return nil, err
	}

	return &i, nil
}

func (i *Interpreter) Toc() string {
	return i.toc
}

func (i *Interpreter) GenerateToc() error {
	var doc io.Reader
	var err error
	if i.document.Path != "" {
		doc, err = os.Open(i.document.Path)
		if err != nil {
			return fmt.Errorf("unable to open document: %w", err)
		}
	} else {
		doc = strings.NewReader(i.document.Content)
	}

	// generate toc tree
	if err := i.parser.Parse(doc); err != nil {
		return fmt.Errorf("unable to parse document: %w", err)
	}

	// build toc string by processing each node
	i.parser.WalkTree(func(header parser.Header, depth int) {
		trimmed := strings.TrimSpace(header.Name)
		var suffix string
		if i.headerCount[trimmed] > 0 {
			suffix = "-" + strconv.Itoa(i.headerCount[trimmed])
		}
		i.headerCount[trimmed]++

		// write table of content line based on header name and depth
		link := strings.Repeat("\t", depth) +
			"* [" + trimmed + "](#" + i.formatHeaderLink(header.Name) +
			suffix + ")\n"

		i.toc += link
	})

	if len(i.toc) > 0 {
		i.toc = "\n" + generatedComment + "\n" + i.toc + "\n" + generatedComment + "\n"
	}

	return nil
}

func (i *Interpreter) formatHeaderLink(text string) string {
	// 1. convert to all lower case
	// 2. convert all dashes to double dashes
	// 3. convert all spaces to dashes
	// 4. remove any non letters to the left
	return strings.TrimLeftFunc(
		strings.Replace(
			strings.Replace(
				strings.ToLower(text),
				"-",
				"--",
				-1,
			),
			" ",
			"-",
			-1,
		),
		func(r rune) bool {
			return !unicode.IsLetter(r)
		})
}

func (i *Interpreter) validate() error {
	var missingDeps []string
	tests := []struct {
		dep string
		chk func() bool
	}{
		{
			dep: "parser",
			chk: func() bool { return i.parser != nil },
		},
	}

	for _, tc := range tests {
		if !tc.chk() {
			missingDeps = append(missingDeps, tc.dep)
		}
	}

	if len(missingDeps) > 0 {
		return fmt.Errorf("unable to initialize interpreter due to (%d) missing depdencies: %s", len(missingDeps), strings.Join(missingDeps, ","))
	}

	return nil
}

type Document struct {
	Path    string
	Content string
}
