# Identicons in Golang

A kinda customizeable identicon implementation based on Bitmessage qidenticon

## Sample icons

![Sample icons](sample.png)

## Example

### Basic

```golang
package main

import (
	"image/png"
	"os"

	"github.com/fivenp/go-identicon"
)

func main() {
	code := qidenticon.Code("test")
	size := 30
	settings := qidenticon.DefaultSettings()
	img := qidenticon.Render(code, size, settings)
	w, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	err = png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}
```

### Advanced

Check the [demo](demo/) folder...
