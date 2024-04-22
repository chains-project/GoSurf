# Arbitrary Code Execution Strategies
The Go language explicitly focuses on [addressing supply chain attacks](https://go.dev/blog/supply-chain). Despite this, in the following, we report techniques that 3rd-party Go dependencies may employ to attain ACE when they are installed or run in the context of downstream project. The categorization comes from [The Hitchhiker's Guide To Malicious Third-Party Dependencies](https://arxiv.org/abs/2307.09087).

For each of these cases, we performed a capability analysis by using Capslock tool to assess the effectiveness of static call graph analysis in identifying capabilities invoked in malicious dependencies when using different injection strategies. 


## Install-Time Execution
There are three types of techniques to achieve ACE when downstream projects install a 3rd-party dependency using package managers, but none of them seems to be applicable to the Go ecosystem.

### [I1] Run command/scripts leveraging install-hooks [Not applicable in Go]
Execution of code by hooking the install process of dependecies in different stages, using specific key-words that package managers may provide.

### [I2] Run code in build script [Not applicable in Go]
Execution of code contained in scripts used by package managers during the installation of dependencies distributed as source code.

### [I3] Run code in build extensions [Not applicable in Go]
Execution of extensions of dependencies that are necessary for their build process. 



## Runtime Execution
There are four techniques to achieve ACE at runtime. The first three seem to be applicable to the Go ecosystem.


### [R1] Insert code in methods executed when importing a module [Applicable in Go]
Attackers can insert malicious code that executes when an import statement is processed, even before the code from the imported module is actually used. In Go, dependencies can execute code upon import in two ways:

*R1.1 By defining an `init()` method. [Examples here](https://itnext.io/golang-stop-trusting-your-dependencies-a4c916533b04).*

```go
package mylib

func init() {
    // malicious code here
}
```

*R1.2 By initializing a global variable with an anonymous function. [Examples here](https://itnext.io/golang-stop-trusting-your-dependencies-a4c916533b04).*

```go
package mylib

var anonym_func string = func() string {
    // malicious code here
    return ""
}()
```

N.B. Importing a package with an underscore prefix prevents Go from automatically removing unused dependencies. This ensures that even if the package is not directly used in the code, its init() function or any anonymous functions assigned to global variables will still be executed.


- **Capslock Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

- **Details**: Identifies the real capability within the method.



### [R2] Insert code in constructors methods [Applicable in Go]
Attackers may target constructor methods as suitable places to insert malicious code because those functions are frequently used in the code to create instances of a struct. While Go doesn't have traditional constructors, developers often define and use common functions as "constructor" functions to initialize structs.


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
    // malicious code here
    return &Person{
        Name: name,
        Age:  age,
    }
}

```

- **Capslock Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

- **Details**: Identifies the real capability within the method.


### [R3] Insert code in commonly-used methods [Applicable in Go]
Attackers may target commonly-used methods within popular imported packages.

For example:

    fmt.Printf() - Formatting output (package: fmt)
    time.Now() - Getting the current time (package: time)
    strconv.Itoa() - Converting an integer to a string (package: strconv)
    http.Get() - Making an HTTP GET request (package: net/http)
    json.Marshal() - Encoding data to JSON (package: encoding/json)
    encoding/base64.StdEncoding.EncodeToString() - Encoding data to base64 (package: encoding/base64)

- **Capslock Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

- **Details**: Identifies the real capability within the method.


### [R4] Run code as build plugin [Not Applicable in Go]
Execute the dependency as a plugin within the build of a downstream project. However, in Go there isn't a direct equivalent for example of Maven plugins for injecting code into the build process. 


## Extend the classification 

### [R5] Run code by using Reflection [Applicable in Go]
Reflection in Go enables dynamic inspection and manipulation of structures, functions, and variables at runtime, facilitating flexible and generic code. Attackers can insert malicious code by exploiting the reflection feature, making challenging to analyze the behavior and the intent of functions and code statically.


```golang
package main

import (
	"fmt"
	"reflect"
)

type Foo struct{}

func (f Foo) Method() {
	// malicious code here
}

func main() {
	var f Foo
	v := reflect.ValueOf(f)
	m := v.MethodByName("Method")
	m.Call(nil) 
}
```

- **Capslock Outcome**: <span style="color:orange">*WEAK FALSE NEGATIVE*</span>

- **Details**: Detects only the `CAPABILITY_REFLECT`, but cannot detect the real capability

- **TODO**: Try other real-world examples, because it might identify the real capability.


### [R6] Run code by using indirect method invocations via interfaces [Applicable in Go]
Attackers can use Go's interface mechanism for dynamic method dispatch. When methods are indirectly invoked via an interface, their specific implementation is determined at runtime, posing challenges for static detection of malicious behavior. 

```golang
package main

type Fooer interface {
    Foo()
}

type Bar struct{}

func (b Bar) Foo() {
    // Malicious code here
}

func invokeFoo(f Fooer) {
    f.Foo()
}

func main() {
    var b Bar
    invokeFoo(b)
}
```

- **Capslock Outcome**: <span style="color:green">*NO FALSE NEGATIVE*</span>

- **Details**: Identifies the real capability within the method.

- **TODO**: Investigate the false positives.


### [R7] Execute imported C code through CGO feature [Applicable in Go]
CGO features enable executing C code in Go binaries. Attackers could exploit this capability to gain more control over the system, and also exploit memory safety concerns related to these low level languages.  

```golang
package main

/*
#include <stdio.h>

void cFunction() {
    // malicious code here
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
- **Capslock Outcome**: <span style="color:orange">*WEAK FALSE NEGATIVE*</span>

- **Details**: Detects only the `CAPABILITY_CGO`, but cannot detect the actual capability.

- **TODO**: Investigate the false positives.

### [R8] Execute dynamically generated code [Applicable in Go]
Attackers can insert functions to dynamically generate and execute code at runtime, creating temporary files, building and executing them. This could make detecting malicious behaviors challenging. 

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
                        // malicious code here
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
    file, _ := os.CreateTemp("", "generated_*.go")
    defer os.Remove(file.Name()) // Clean up temporary file

    // Write the generated code to the temporary file
    file.WriteString(code);
    file.Close();

    // Run the Go file
    cmd := exec.Command("go", "run", file.Name())
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}

```

- **Capslock Outocome**: <span style="color:orange">*WEAK FALSE NEGATIVE*</span>

- **Details**: Detects only the `CAPABILITY_EXEC`, but cannot detect the actual capabilities.


### [R9] Execute pre-built code loaded at runtime [Applicable in Go]
Attackers could load pre-built code at runtime making the detection of malicious behavior challenging. In Go, this is possibile mainly in two ways:

*R9.1 By importing and using external binary `plugins`*

```go
// plugin.go
package main

func PluginFunc() {
	// malicious code here
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
	p, _ := plugin.Open("./plugin.so")
	
    	// Look up the symbol (function) from the loaded plugin
	sym, _ := p.Lookup("PluginFunc")
	
    	// Assert and call the function if found
	if fn, ok := sym.(func()); ok {
		fn()
	} else {
		fmt.Println("PluginFunc has unexpected type")
	}
}
```

*R9.2 By using the `os.exec` package to execute arbitrary external commands.*

```go
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    binaryPath := "/path/to/prebuilt/binary"
    cmd := exec.Command(binaryPath)
    _ := cmd.Run()
}
```

- **Outcomes**: <span style="color:orange">*WEAK FALSE NEGATIVE*</span>

- **Details**: Detects only the `CAPABILITY_EXEC`, but cannot detect the actual capabilities.


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

