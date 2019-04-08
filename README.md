# request

## Golang HTTP Requests for Humans™ ✨🍰✨

### Usage

* 基本用法

```(golang)

requests.Get("https://golang.org)
requests.Post("https://golang.org, "application/json", `{"a": "b"}`)
```

* 高级用法

```(golang)
package main

import (
    "log"

    "github.com/luopengift/requests"
)

func main() {
    sess := requests.New()                                           // 创建session
    sess.SetHeader("aaa", "bbb")                                     // 全局配置, 会追加到使用这个sess的所有请求中
    req, err := requests.NewRequest("GET", "http://httpbin.org", nil) // 创建一个GET请求
    if err != nil {
        log.Fatal(err)
        return
    }
    req.SetHeader("foo", "bar") // req的参数设置会覆盖sess中的参数
    resp, err := sess.DoRequest(req) //发送创建的请求
    if err != nil {
        log.Fatal(err)
        return
    }
    _, err = resp.Text() //解析响应
    if err != nil {
        log.Fatal(err)
        return
    }
}
```
