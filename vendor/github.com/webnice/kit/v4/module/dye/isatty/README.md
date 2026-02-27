# isatty

isatty for golang

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/webnice/kit/v4/module/dye/isatty"
)

func main() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("Is Terminal")
	} else if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		fmt.Println("Is Cygwin/MSYS2 Terminal")
	} else {
		fmt.Println("Is Not Terminal")
	}
}
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
