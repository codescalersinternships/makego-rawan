# Makego

This repository provides a Go-based tool to execute Makefile targets with support for concurrency.

## Features

-   Execute Makefile targets directly from Go.
-   Concurrency support: independent stages can run in parallel.
-   Ordered execution for dependent stages.
-   Error handling through channels for safe concurrent execution.


## Installation

1.  Clone this repository:

    ``` bash
    git clone https://github.com/codescalersinternships/makego-rawan.git
    cd makego-rawan
    ```

## Usage

### Code Example 

```go
import makego "github.com/codescalersinternships/makego-rawan/pkg"

func main() {
	err := makego.ExecuteMakefile(path, targets)
	if err != nil {
		panic(err)
	}
}
```

### Running a Makefile Target

``` bash
go run cmd/main.go -f "testdata/makefile" target
```

### Example

Given the following Makefile:

``` makefile
hello: hello1.txt hello2.txt 
	cat hello1.txt
	cat hello2.txt

hello1.txt: 
	echo "Hello 1, Make!" > hello1.txt

hello2.txt: 
	echo "Hello 2, Make!" > hello2.txt

clean:
	rm -f hello1.txt
	rm -f hello2.txt

```

Running:

``` bash
go run cmd/main.go -f "testdata/makefile" hello
```

Will produce:

    echo "Hello 2, Make!" > hello2.txt
    echo "Hello 1, Make!" > hello1.txt
    cat hello1.txt
    Hello 1, Make!
    cat hello2.txt
    Hello 2, Make!

### Flags

- `-f`: The makefile to be executed (must be provided)

### Arguments

- `<target>` : The target to be run (optional)
    - default: first target in the makefile

## Concurrency

-   Independent targets run in parallel
-   Channels are used to handle errors