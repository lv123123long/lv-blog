package main

import (
	"encoding/json"
	"fmt"
)

// 定义泛型结构体 Response
type Response[T any] struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    T        `json:"data"`
}

func main() {
	// 创建一个包含字符串数据的 Response 实例
	strResponse := Response[string]{
		Code:    200,
		Message: "Success",
		Data:    "Hello, World!",
	}

	// 将 Response 实例转换为 JSON 字符串
	jsonStr, err := json.Marshal(strResponse)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonStr))

	// 创建一个包含整数数据的 Response 实例
	intResponse := Response[int]{
		Code:    400,
		Message: "Error",
		Data:    12345,
	}

	// 将 Response 实例转换为 JSON 字符串
	jsonStr, err = json.Marshal(intResponse)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonStr))
}
