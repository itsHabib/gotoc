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

// Interpreter is responsible for generating the table of contents string from
// the parsers tree.
type Interpreter struct {
	document    Document
	headerCount map[string]int
	indent      int
	parser      *parser.Parser
	toc         string
}

// New returns an instantiated instance of an Interpreter. An interpreter
// expects a document that contains either a path or content.
func New(document Document) (*Interpreter, error) {
	if document.Path != "" {
		if _, err := os.Stat(document.Path); os.IsNotExist(err) {
			return nil, errors.New("unable to create interpreter, no file found at document path")
		}
	} else if document.Content == "" {
		return nil, errors.New("unable to create interpreter, content is empty")
	}

	return &Interpreter{
		parser:      parser.New(),
		headerCount: make(map[string]int),
		document:    document,
	}, nil
}

// GenerateToc will use the parser to generate a tree from the document. Once
// constructed the interpreter will walk the tree and create the header string
// based on the tree node string and depth.
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

	// build toc string by processing each node, use DFS to maintain proper
	// order in a toc
	i.parser.WalkTree(func(header string, depth int) {
		trimmed := strings.TrimSpace(header)
		var suffix string
		if i.headerCount[trimmed] > 0 {
			suffix = "-" + strconv.Itoa(i.headerCount[trimmed])
		}
		i.headerCount[trimmed]++

		// write table of content line based on header name and depth
		i.toc += strings.Repeat("\t", depth) +
			"* [" + trimmed + "](#" + i.formatHeaderLink(trimmed) +
			suffix + ")\n"
	})

	if len(i.toc) > 0 {
		i.toc = generatedComment + "\n" + i.toc + generatedComment + "\n"
	}

	return nil
}

// Toc returns the table of contents string
func (i *Interpreter) Toc() string {
	return i.toc
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
