package main

import "os"
import "encoding/json"

func main() {
	// CAPABILITY_REFLECT, 1°CAPABILITY_READ_SYSTEM_STATE
	json.Marshal(os.Environ())      
	// 2°CAPABILTIY_READ_SYSTEM_STATE 
	os.Getenv("example")
}
