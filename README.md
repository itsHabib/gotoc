# GoToC - Markdown ToC generator

## Build

```bash
go build -o gotoc cmd/main.go
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