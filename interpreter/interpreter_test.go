package interpreter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateToc(t *testing.T) {
	for _, tc := range []struct {
		description string
		document    Document
		expectedToc string
		wantErr     bool
	}{
		{
			description: "Multi Branch tree",
			document: Document{
				Content: `
# Header 1
hello this is header one

## Header 2
### Header 3
test tester test

## Header 2 2

## Header 2 3

### Header 3 3

#### Header 4 3

###### Header 6 3

### Header 3 4

## Header 2 4

### Header 3 4
`,
			},
			expectedToc: "\n" + generatedComment + "\n" + `* [Header 1](#header-1)
	* [Header 2](#header-2)
		* [Header 3](#header-3)
	* [Header 2 2](#header-2-2)
	* [Header 2 3](#header-2-3)
		* [Header 3 3](#header-3-3)
			* [Header 4 3](#header-4-3)
				* [Header 6 3](#header-6-3)
		* [Header 3 4](#header-3-4)
	* [Header 2 4](#header-2-4)
		* [Header 3 4](#header-3-4-1)
` + "\n" + generatedComment + "\n",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			interpreter, err := NewInterpreter(tc.document)
			require.NoError(t, err)

			err = interpreter.GenerateToc()
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.expectedToc, interpreter.toc)
		})

	}
}
