package main

import (
	"encoding/json"
	"log"
	"os"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

func main() {

	jsonSource, err := os.ReadFile("data.json")
	if err != nil {
		panic(err)
	}

	jsonData := make(map[string]interface{})

	err = json.Unmarshal(jsonSource, &jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft4

	sch, err := compiler.Compile("user.json")
	if err != nil {
		panic(err)
	}

	err = sch.Validate(jsonData)
	if err != nil {
		log.Println(err)
	}

}
