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
			description: "Invalid Double header 1",
			document: `
# header 1
## header 2
# header 1 again invalid
`,
			wantErr: true,
		},
		{
			description: "Mixed header depths",
			document: `
## header 2
##### header 5
###### header 6
#### header 4
### header 3
##### header 5
## header 2 2
`,
			expectedTree: &tree{
				root: &node{
					children: []*node{
						{
							depth: 1,
							header: header{
								name: " header 2",
								text: "## header 2",
							},
							children: []*node{
								{
									depth: 4,
									header: header{
										name: " header 5",
										text: "##### header 5",
									},
									children: []*node{
										{
											depth: 5,
											header: header{
												name: " header 6",
												text: "###### header 6",
											},
										},
									},
								},
								{
									depth: 3,
									header: header{
										name: " header 4",
										text: "#### header 4",
									},
								},
								{
									depth: 2,
									header: header{
										name: " header 3",
										text: "### header 3",
									},
									children: []*node{
										{
											depth: 4,
											header: header{
												name: " header 5",
												text: "##### header 5",
											},
										},
									},
								},
							},
						},
						{
							depth: 1,
							header: header{
								name: " header 2 2",
								text: "## header 2 2",
							},
						},
					},
				},
			},
		},
		{
			description: "Multi branched tree",
			document: `
# header 1
hello this is header one

## header 2
### header 3
test tester test

## header 2 2

## header 2 3

### header 3 3

#### header 4 3

###### header 6 3

### header 3 4

## header 2 4

### header 3 4
`,
			expectedTree: &tree{
				root: &node{
					depth: 0,
					header: header{
						name: " header 1",
						text: "# header 1",
					},
					children: []*node{
						{
							depth: 1,
							header: header{
								name: " header 2",
								text: "## header 2",
							},
							children: []*node{
								{
									depth: 2,
									header: header{
										name: " header 3",
										text: "### header 3",
									},
								},
							},
						},
						{
							depth: 1,
							header: header{
								name: " header 2 2",
								text: "## header 2 2",
							},
						},
						{
							depth: 1,
							header: header{
								name: " header 2 3",
								text: "## header 2 3",
							},
							children: []*node{
								{
									depth: 2,
									header: header{
										name: " header 3 3",
										text: "### header 3 3",
									},
									children: []*node{
										{
											depth: 3,
											header: header{
												name: " header 4 3",
												text: "#### header 4 3",
											},
											children: []*node{
												{
													depth: 5,
													header: header{
														name: " header 6 3",
														text: "###### header 6 3",
													},
												},
											},
										},
									},
								},
								{
									depth: 2,
									header: header{
										name: " header 3 4",
										text: "### header 3 4",
									},
								},
							},
						},
						{
							depth: 1,
							header: header{
								name: " header 2 4",
								text: "## header 2 4",
							},
							children: []*node{
								{
									depth: 2,
									header: header{
										name: " header 3 4",
										text: "### header 3 4",
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
# header 1
hello this is header one

## header 2
### header 3
test tester test

#### header 4
something herre
##### header 5
something over there
###### header 6
last header to use

####### header 7 this should be ignored
ignore this header above
`,
			expectedTree: &tree{
				root: &node{
					depth: 0,
					header: header{
						name: " header 1",
						text: "# header 1",
					},

					children: []*node{
						{
							depth: 1,
							header: header{
								name: " header 2",
								text: "## header 2",
							},
							children: []*node{
								{
									depth: 2,
									header: header{
										name: " header 3",
										text: "### header 3",
									},
									children: []*node{
										{
											depth: 3,
											header: header{
												name: " header 4",
												text: "#### header 4",
											},
											children: []*node{
												{
													depth: 4,
													header: header{
														name: " header 5",
														text: "##### header 5",
													},
													children: []*node{
														{
															depth: 5,
															header: header{
																name: " header 6",
																text: "###### header 6",
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
			parser := New()
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
