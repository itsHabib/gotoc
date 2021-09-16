package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		description  string
		document     string
		expectedTree *tree
		wantErr      bool
	}{
		{
			description: "No headers in document",
			document: `
no headers at all
in this whole doc
this will produce a nil tree
`,
		},
		{
			description: "Invalid Double Header 1",
			document: `
# Header 1
## Header 2
# Header 1 again invalid
`,
			wantErr: true,
		},
		{
			description: "Mixed header depths",
			document: `
## Header 2
##### Header 5
###### Header 6
#### Header 4
### Header 3
##### Header 5
## Header 2 2
`,
			expectedTree: &tree{
				root: &node{
					children: []*node{
						{
							depth: 1,
							header: Header{
								Name: " Header 2",
								Text: "## Header 2",
							},
							children: []*node{
								{
									depth: 4,
									header: Header{
										Name: " Header 5",
										Text: "##### Header 5",
									},
									children: []*node{
										{
											depth: 5,
											header: Header{
												Name: " Header 6",
												Text: "###### Header 6",
											},
										},
									},
								},
								{
									depth: 3,
									header: Header{
										Name: " Header 4",
										Text: "#### Header 4",
									},
								},
								{
									depth: 2,
									header: Header{
										Name: " Header 3",
										Text: "### Header 3",
									},
									children: []*node{
										{
											depth: 4,
											header: Header{
												Name: " Header 5",
												Text: "##### Header 5",
											},
										},
									},
								},
							},
						},
						{
							depth: 1,
							header: Header{
								Name: " Header 2 2",
								Text: "## Header 2 2",
							},
						},
					},
				},
			},
		},
		{
			description: "Multi branched tree",
			document: `
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
			expectedTree: &tree{
				root: &node{
					depth: 0,
					header: Header{
						Name: " Header 1",
						Text: "# Header 1",
					},
					children: []*node{
						{
							depth: 1,
							header: Header{
								Name: " Header 2",
								Text: "## Header 2",
							},
							children: []*node{
								{
									depth: 2,
									header: Header{
										Name: " Header 3",
										Text: "### Header 3",
									},
								},
							},
						},
						{
							depth: 1,
							header: Header{
								Name: " Header 2 2",
								Text: "## Header 2 2",
							},
						},
						{
							depth: 1,
							header: Header{
								Name: " Header 2 3",
								Text: "## Header 2 3",
							},
							children: []*node{
								{
									depth: 2,
									header: Header{
										Name: " Header 3 3",
										Text: "### Header 3 3",
									},
									children: []*node{
										{
											depth: 3,
											header: Header{
												Name: " Header 4 3",
												Text: "#### Header 4 3",
											},
											children: []*node{
												{
													depth: 5,
													header: Header{
														Name: " Header 6 3",
														Text: "###### Header 6 3",
													},
												},
											},
										},
									},
								},
								{
									depth: 2,
									header: Header{
										Name: " Header 3 4",
										Text: "### Header 3 4",
									},
								},
							},
						},
						{
							depth: 1,
							header: Header{
								Name: " Header 2 4",
								Text: "## Header 2 4",
							},
							children: []*node{
								{
									depth: 2,
									header: Header{
										Name: " Header 3 4",
										Text: "### Header 3 4",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			description: "Single list tree with all depths covered",
			document: `
# Header 1
hello this is header one

## Header 2
### Header 3
test tester test

#### Header 4
something herre
##### Header 5
something over there
###### Header 6
last header to use

####### Header 7 this should be ignored
ignore this header above
`,
			expectedTree: &tree{
				root: &node{
					depth: 0,
					header: Header{
						Name: " Header 1",
						Text: "# Header 1",
					},

					children: []*node{
						{
							depth: 1,
							header: Header{
								Name: " Header 2",
								Text: "## Header 2",
							},
							children: []*node{
								{
									depth: 2,
									header: Header{
										Name: " Header 3",
										Text: "### Header 3",
									},
									children: []*node{
										{
											depth: 3,
											header: Header{
												Name: " Header 4",
												Text: "#### Header 4",
											},
											children: []*node{
												{
													depth: 4,
													header: Header{
														Name: " Header 5",
														Text: "##### Header 5",
													},
													children: []*node{
														{
															depth: 5,
															header: Header{
																Name: " Header 6",
																Text: "###### Header 6",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			parser := NewParser()
			doc := strings.NewReader(tc.document)

			err := parser.Parse(doc)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.expectedTree, parser.tree)
		})
	}
}
