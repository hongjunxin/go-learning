package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 定义一个 map
	data := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
		"asset": map[string]interface{}{
			"house": 1,
			"car":   2,
		},
	}

	// 将 map 转换为 json 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	fmt.Println("JSON data:", string(jsonData))

	// 将 json 字符串转换为 map
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("json.Unmarshal error:", err)
		return
	}

	fmt.Println("Result map:", result)
	fmt.Printf("age: %v, house %v\n", result["age"], result["asset"].(map[string]interface{})["house"])
}
