package main

import (
	"fmt"
	"github.com/alanhi/sfsdk"
)

func main() {
	client := sfsdk.NewClient("YOUR_VALUE", "YOUR_VALUE", sfsdk.Test)

	res, err := client.Execute("YOUR_VALUE", "{}")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.IsSuccess())
	fmt.Println(res.String())
	fmt.Println(res.Json())
}
