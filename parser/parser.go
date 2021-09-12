package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	markDownMaxHeaderDepth = 6
	markDownHeaderChar     = '#'
)

type Parser struct {
	depth    int
	cursor   int
	path     []*node
	tree     *tree
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(doc io.Reader) error {
	scanner := bufio.NewScanner(doc)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}

		// consume all '#' symbols to get p.depth of current header
		p.consumeHeader(text)

		// no header found
		if p.cursor == 0 {
			continue
		}

		if p.cursor == 1 && p.depth != 0 {
			return errors.New("unable to parse markdown document with more than one main header '#'. ")
		}

		n := node{
			depth:  p.cursor - 1,
			header: text[p.cursor:],
			text:   text,
		}

		// root, this locks the root depth at p.cursor -1. If we encounter a
		// depth that is higher, the parser will throw an error
		if p.depth == 0 && p.cursor > 0 {
			p.initTree(&n)
			continue
		}

		// new child header
		if p.cursor-1 >= p.depth && p.cursor <= markDownMaxHeaderDepth {
			if err := p.newChild(&n); err != nil {
				return fmt.Errorf("unable to add new child: %w", err)
			}
			continue
		}

		// encountered a header with equal or higher depth
		if p.cursor-1 < p.depth {
			// encountered a depth that is higher than the root depth, invalid
			// document
			if p.cursor-1 < p.tree.root.depth {
				return fmt.Errorf("can not climb path past the root depth: %d", p.cursor-1)
			}

			// climb up the path until we get to the nearest node in depth that
			// can be a parent to the current node. Set depth to valid
			// parent depths
			if err := p.climbPath(p.cursor - 1); err != nil {
				return fmt.Errorf("unable to climb path: %w", err)
			}

			// add new child to correct parent
			if err := p.newChild(&n); err != nil {
				return fmt.Errorf("unable to add new child: %w", err)
			}
		}
	}

	return nil
}

func (p *Parser) climbPath(depth int) error {
	if depth < p.tree.root.depth {
		return fmt.Errorf("can not climb path past the root depth: %d, given: %d", p.tree.root.depth, depth)
	}

	// pop items in path until we reach a node with the needed depth
	for len(p.path) > 0 {
		cur := p.path[len(p.path)-1]
		if cur.depth < depth {
			p.depth = cur.depth
			return nil
		}
		p.path = p.path[:len(p.path)-1]
	}

	return errors.New("unable to find needed parent depth")
}

// climb up the path until we get to the correct parent
func (p *Parser) consumeHeader(line string) {
	p.cursor = 0
	for ; p.cursor < len(line) && line[p.cursor] == markDownHeaderChar; p.cursor++ {
	}
}

func (p *Parser) initTree(n *node) {
	// if we never found a main header #, we still need to create a dummy root
	// node to hold all children nodes
	if n.depth != 0 {
		p.tree = &tree{
			root: &node{
				depth:    0,
				children: []*node{n},
			},
		}
		p.path = append(p.path, p.tree.root)
	} else {
		p.tree = &tree{n}
	}

	p.depth = p.cursor
	p.path = append(p.path, n)
}

func (p *Parser) newChild(n *node) error {
	if len(p.path) == 0 {
		return errors.New("unable to add a new child with no parent in the path")
	}

	current := p.path[len(p.path)-1]
	current.children = append(current.children, n)
	p.path = append(p.path, n)
	p.depth = p.cursor

	return nil
}

type tree struct {
	root *node
}

type node struct {
	children []*node
	depth    int
	text     string
	header   string
}
