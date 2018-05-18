package main

import (
	json "github.com/json-iterator/go"
	"fmt"
)

type JsonStruct struct {
	Name string
	Age  int
}

func main() {
	jsonDemo := &JsonStruct{
		Name: "dubing",
		Age:  30,
	}
	data, err := json.Marshal(jsonDemo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))

	json1 := new(JsonStruct)
	fmt.Println("UnMarshall json1")
	err = json.Unmarshal(data, json1)
	fmt.Println(fmt.Sprintf("%+v", json1))
}
