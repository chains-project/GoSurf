// main.go
package main

import (
	"fmt"
	"plugin"
)

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
