package main

import "github.com/orvyl/kit/id"
import "fmt"

func main() {
	idGen, err := id.NewGenerator(id.Settings{IsAlphaNumeric: true})
  if err != nil {
    fmt.Printf("ERR %v", err)
  }
  id, err := idGen.Next()

	fmt.Printf("%v %v", id, err)
}
