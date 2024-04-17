# Arbitrary Code Execution Strategies
The Go language explicitly focuses on [addressing supply chain attacks](https://go.dev/blog/supply-chain). Despite this, in the following, we report techniques that 3rd-party Go dependencies may employ to attain ACE when they are installed or run in the context of downstream project. The categorization comes from [The Hitchhiker's Guide To Malicious Third-Party Dependencies](https://arxiv.org/abs/2307.09087).

For each of these cases, we performed a capability analysis by using Capslock tool to assess the effectiveness of static call graph analysis in identifying capabilities invoked in malicious dependencies when using different injection strategies. 


## Install-Time Execution
There are three types of techniques to achieve ACE when downstream projects install a 3rd-party dependency using package managers, but none of them seems to be applicable to the Go ecosystem.

### [I1] Run command/scripts leveraging install-hooks [Not applicable in Go
Execution of code by hooking the install process of dependecies in different stages, using specific key-words that package managers may provide.

### [I2] Run code in build script [Not applicable in Go]
Execution of code contained in scripts used by package managers during the installation of dependencies distributed as source code.

### [I3] Run code in build extensions [Not applicable in Go]
Execution of extensions of dependencies that are necessary for their build process. 



## Runtime Execution
There are four scenarios where malicious code can be executed at runtime. The first three seem to be applicable to the Go ecosystem.


### [R1] Insert code in methods executed when importing a module [Applicable in Go]
    
Execution of code when an import statement is processed, even before the code from the imported module is actually used. In Go, dependencies can execute code upon import in two ways:

1. By defining an `init()` method. [Examples here](https://itnext.io/golang-stop-trusting-your-dependencies-a4c916533b04).

    ```golang
    package mylib

    func init() {
        // malicious code here
    }
    ```
    If the evilpkg is imported with an underscore prefix, which prevents Go's automatic removal of unused dependencies, this ensures that event tough the package is not directly used in the code, its init() function wil till be executed.


2. By initializing a variable with an anonymous function. [Examples here](https://itnext.io/golang-stop-trusting-your-dependencies-a4c916533b04).

    ```go
    package mylib

    var anonym_func string = func() string {
        // malicious code here
        return ""
    }()
    ```
    When this package is imported into another Go program, this initialization code will execute.


**Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

**Details**: Identifies the real capability within the method.



### [R2] Insert code in constructors methods [Applicable in Go]
Attackers may target constructors methods as suitable places to insert malicious code. 

Go have a mechanism called "struct initialization" which can be considered somewhat similar to constructor methods. In Go, structs are used to define types with a collection of fields. When a struct is initialized, all its fields are initialized to their zero values by default.

An attacker could potentially insert malicious code into a function that initializes a struct or into a function that is commonly used to create instances of a struct. 

```go                                                 
package mylib
import "fmt"

// Define a struct type
type Person struct {
    Name string
    Age  int
}

// Constructor-like function for Person
func NewPerson(name string, age int) *Person {
    fmt.Println("malicious code here")
    return &Person{
        Name: name,
        Age:  age,
    }
}

```

**Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

**Details**: Identifies the real capability within the method.


### [R3] Insert code in commonly-used methods [Applicable in Go]
Attackers may target methods within 3rd-party dependencies that they distribute to downstream users. These methods are attractive because they are commonly used, increasing the likelihood that downstream users will invoke and rely on them.


### [R4] Run code as build plugin [Not Applicable in Go]
Execute the dependency as a plugin within the build of a downstream project. However, in Go there isn't a direct equivalent of Maven plugins for injecting code into the build process. 


## Extend the classification 

### [R5] Indirect method invocations via interfaces [Applicable in Go]
Go's interface mechanism allows for dynamic method dispatch. If a method is invoked indirectly via an interface, the specific method implementation is determined at runtime.

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


### [R6] Execute imported C code through cgo feature [Applicable in Go]
Execute code from different languages in Go binaries. Go packages can interact with C code using cgo feature. 

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

### [R7] Execute dynamically generated code [Applicable in Go]
Execute code generated at runtime. Go packages can generate functions or other code at runtime by creating temporary files, building them and executing them.

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


### [R8] Execute dynamically loaded code [Applicable in Go]
Go packages can loads at runtime code (e.g., via plugins or external modules).

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
#
### Use of Reflection

Call graph analysis typically examines the code at rest, without executing it. If your Go package uses reflection or other forms of dynamic code execution to invoke methods, these may not be captured by static analysis.

```golang
package main

import (
	"fmt"
	"reflect"
)

type Foo struct{}

func (f Foo) Method() {
	// Capability invocation here
}

func main() {
	var f Foo
	v := reflect.ValueOf(f)
	m := v.MethodByName("Method")
	m.Call(nil) 
}
```

**Outcome**: WEAK FALSE NEGATIVE.

**Details**: Detects only the CAPABILITY_REFLECT, but cannot detect the real capability

**TODO**: Try other real-world examples, because it might identify the real capability.
