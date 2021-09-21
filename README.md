# GoToC - Markdown ToC generator

## Build

```bash
go build -o gotoc cmd/main.go
```

## Usage
```bash
Usage of ./gotoc:
  -file string
    	path to a file to generate the table of contents from.
  -print
    	only print the table of contents string and not write it to file. the table of contents are printed by default, unless -write is used.
  -text string
    	generate a table of contents from an input string. Can not be used in conjunction with the -file flag.
  -write
    	write the table of contents string to given file. gotoc will attempt to write the contents after the root header #. If no root header is found the contents are written at the end of the file. This can not be used in conjunction with -print.
```

## Examples

### Printing Toc
```bash
./gotoc -file ./cmd/EXAMPLE.md
<!-- gotoc generated table of contents -->
* [Tester](#tester)
	* [Tester 2](#tester-2)
		* [Test 3](#test-3)
	* [Tester 2 2](#tester-2-2)
		* [Tester 3 2](#tester-3-2)
			* [Tester 4](#tester-4)
```

### Writing Toc
```bash
./gotoc -file ./cmd/EXAMPLE.md -write
```

### Using Text
```bash
./gotoc -text "# Hello there \

## hi "
<!-- gotoc generated table of contents -->
* [Hello there](#hello-there-)
	* [hi](#hi-)
<!-- gotoc generated table of contents -->

```