# Limitations in static capability analysis

Call graph analysis is a powerful tool for understanding capabiltities within a Go package, but there are some cases where it may generate false negatives, meaning it fails to identify all method invocations. Here are some scenarios to consider.


#
### Interoperability with cgo 
If your Go package interacts with C code using cgo, these method invocations may not be accurately represented in the call graph.

```golang
package main

/*
#include <stdio.h>

void cFunction() {
	// Capability invocation here
}
*/
import "C"

type Invoker struct{}

func (i Invoker) InvokeMethod() {
	C.cFunction()
}

func main() {
	invoker := Invoker{}
	invoker.InvokeMethod()
}
```

**Outcome**: <span style="color:orange">*WEAK FALSE NEGATIVE*</span>

**Details**: Detects only the `CAPABILITY_CGO`, but cannot detect the actual capability.

**TODO**: Investigate the false positives.

#
### Method invocations in dynamically generated code
If your Go package uses code dynamically generated, the generated code may not be included in the call graph analysis.

```golang
package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Foo struct{}

func (f Foo) Method() {}

func GenerateCode() string {
		return `package main

				type Foo struct{}

				func (f Foo) Method() {
						// Capability invocation here
				}

				func main() {
						f := Foo{}
						f.Method()
				}
				`
}

func main() {
	code := GenerateCode()

	// Create a temporary Go file to hold the generated code
	file, err := os.CreateTemp("", "generated_*.go")
	if err != nil {
			fmt.Println("Error creating temporary file:", err)
			return
	}
	defer os.Remove(file.Name()) // Clean up temporary file

	// Write the generated code to the temporary file
	if _, err := file.WriteString(code); err != nil {
			fmt.Println("Error writing to temporary file:", err)
			return
	}

	// Close the file
	if err := file.Close(); err != nil {
			fmt.Println("Error closing temporary file:", err)
			return
	}

	// Run the Go file as a subprocess
	cmd := exec.Command("go", "run", file.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
			fmt.Println("Error executing generated code:", err)
			return
	}
}

```

**Outocome**: <span style="color:red">*STRONG FALSE NEGATIVE*</span>

**Details**: It does not identify any capability in the dynamically generated code.

**TODO**: Quantify these cases in real-world packages.

#
### Dynamic Code Loading
If your go package loads dynamically codes at runtime (e.g., via plugins or external modules), static analysis may not be sufficient to trace the function invocations.

```go
// plugin.go
package main

func PluginFunc() {
	// Capability Inovcation here
}
```

```bash
go build -buildmode=plugin -o plugin.so plugin.go
```

```go
package main

import "fmt"
import "plugin"

func main() {
	// Load the plugin dynamically
	p, err := plugin.Open("./plugin.so")
	if err != nil {
		fmt.Println("Error loading plugin:", err)
		return
	}

	// Look up the symbol (function) from the loaded plugin
	sym, err := p.Lookup("PluginFunc")
	if err != nil {
		fmt.Println("Error looking up symbol:", err)
		return
	}

	// Assert and call the function if found
	if fn, ok := sym.(func()); ok {
		fn()
	} else {
		fmt.Println("PluginFunc has unexpected type")
	}
}
```

**Outcomes**: <span style="color:red">*STRONG FALSE NEGATIVE*</span>

**Details**: It does not identify any capability in the (pre-compiled) dynamically imported plugin.

**TODO**: Quantify these cases in real-world packages.


#
### Indirect method invocations via interfaces
Go's interface mechanism allows for dynamic method dispatch. If a method is invoked indirectly via an interface, and the specific method implementation is determined at runtime, this may not be accurately captured in the call graph.

```golang
package main

type Fooer interface {
	Foo()
}

type Bar struct{}

func (b Bar) Foo() {
	// Capability invocation here
}

func invokeFoo(f Fooer) {
	f.Foo()
}

func main() {
	var b Bar
	invokeFoo(b)
}
```

**Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

**Details**: Identifies the real capability within the method (e.g., `CAPABILITY_NETWORK`, `CAPABILITY_EXEC`).

**TODO**: Investigate the false positives.

#
### Buggy Case to analyze

The second occurence of `CAPABILITY_READ_SYSTEM_STATE` is not identified

```go     
package main

import "os"
import "encoding/json"

func main() {
	// CAPABILITY_REFLECT, 1°CAPABILITY_READ_SYSTEM_STATE
	json.Marshal(os.Environ())      
	// 2°CAPABILTIY_READ_SYSTEM_STATE 
	os.Getenv("example")
}
```

In the following case, the `CAPABILITY_READ_SYSTEM_STATE` is identified
```go     
package main

import "os"

func main() {
	// CAPABILTIY_READ_SYSTEM_STATE 
	os.Getenv("example")
}
```

