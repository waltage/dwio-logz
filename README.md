# DWIO logz

```go
package main

import (
	"fmt"
	"os"

	"github.com/waltage/dwio/logz"
)

var log = logz.DefaultLog("main")
var log2 = logz.NewLog(os.Stderr, "packageName", logz.Info)

func main() {
	log.Level = logz.Debug

	log.Debug("Testing a debug")

	log.Warn("A warning")
	log.WarnMsg(map[string]string{"testing": "yes", "testing2": "no"},
		"Warning with JSON Struct")

	fmt.Println("\n...Reset Level...\n")
	log.Level = logz.Warn

	log.Debug("Debugging message no longer prints")

	log2.Info("Messag from second logger")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Errors are recoverable!")
			fmt.Println("-->", r)
		}

		log2.Fatal("Fatals are not...")
	}()

	log2.Error("Panic!")
}
```

Results in:
```
2021/01/14 18:43:12.696270 [DEBUG] [main] main.go:16: Testing a debug
2021/01/14 18:43:12.696364 [WARN ] [main] main.go:18: A warning
2021/01/14 18:43:12.696434 [WARN ] [main] main.go:19: Warning with JSON Struct | <msg>{"testing":"yes","testing2":"no"}</msg>

...Reset Level...

2021/01/14 18:43:12.696440 [INFO ] [packageName] main.go:27: Messag from second logger
2021/01/14 18:43:12.696443 [ERROR] [packageName] main.go:38: Panic!
Errors are recoverable!
--> main.go:38: Panic!
2021/01/14 18:43:12.696465 [FATAL] [packageName] main.go:35: Fatals are not...
```