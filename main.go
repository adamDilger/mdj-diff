package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("pbt.mdj")
	if err != nil {
		panic(err)
	}

	var root Project
	json.NewDecoder(file).Decode(&root)

	fmt.Println("mdj-diff", root.OwnedElements[0])
}
