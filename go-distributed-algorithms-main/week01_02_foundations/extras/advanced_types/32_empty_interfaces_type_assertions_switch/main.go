package main

import "fmt"

func printValue(value any) {
	fmt.Println(value)
}

func checkType(value interface{}) {
	switch v := value.(type) {
	case int:
		fmt.Println("Integer:", v)
	case string:
		fmt.Println("String:", v)
	case float64:
		fmt.Println("Float64:", v)
	case bool:
		fmt.Println("Boolean:", v)
	default:
		fmt.Println("Unknown type")
	}
}

func main() {
	mixedSlice := []interface{}{42, "hello", 3.14, true}

	for _, v := range mixedSlice {
		checkType(v)
	}

	var emptyInterface any
	emptyInterface = "Hello World!"

	if str, ok := emptyInterface.(string); ok {
		fmt.Println("String value:", str)
	} else {
		fmt.Println("Not a string")
	}

}
