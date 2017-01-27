# qidenticon (from Bitmessage) ported to Go



## Sample icons

![Sample icons](sample.png)

## Example

```golang
package main

import (
	"image/png"
	"os"

	"github.com/jakobvarmose/go-qidenticon"
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
