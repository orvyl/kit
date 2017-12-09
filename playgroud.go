package main

import "github.com/orvyl/kit/id"
import "fmt"

func main() {
	idGen, err := id.NewGenerator(true, id.Settings{UseAWSData: false})
	if err != nil {
		fmt.Printf("ERR %v", err)
	}
	id, err := idGen.Next()

	fmt.Printf("%v %v", id, err)
}
