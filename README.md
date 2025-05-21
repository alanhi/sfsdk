### SF SDK

顺丰丰桥sdk，实现了基础鉴权逻辑

使用示例：

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/alanhi/sfsdk"
)

func main() {
	client := sfsdk.NewClient("_YOUR_VALUE_", "_YOUR_VALUE_", sfsdk.Test)

	body := make(map[string]any, 0)
	body["trackingType"] = "1"
	body["trackingNumber"] = []string{"SF123456789"}
	body["checkPhoneNo"] = "1234"

	bodyBytes, _ := json.Marshal(body)
	bodyJson := string(bodyBytes)
	fmt.Println("req:", bodyJson)

	apiRes, err := client.Execute("EXP_RECE_SEARCH_ROUTES", bodyJson)
	fmt.Println("res:", apiRes.Json(), err)
}
```
