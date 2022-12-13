a package for recording system metrics on linux.

`go get github.com/rexlx/performance`

```go
package main

import (
	"fmt"

	"github.com/rexlx/performance"
)

func main() {
	stream := make(chan []*performance.DiskStat)
	go performance.GetDiskUsage(stream, 1)
	msg := <-stream
	for _, i := range msg {
		fmt.Println(*i)
	}

}
```
