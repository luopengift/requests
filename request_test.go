package requests_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/luopengift/requests"
)

func Test_Get(t *testing.T) {
	t.Log("Testing get request")
	sess := requests.New().SetRetry(3)
	sess.SetHeader("a", "b").SetCookie("username", "golang").SetAuth("user", "123456")
	req := requests.NewRequest("GET", "http://httpbin.org/get", nil)
	req.SetParam("uid", 1).SetCookie("username", "000000")
	resp, err := sess.DoRequest(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Log(resp.StatusCode, err)
	t.Log(resp.Text())
	// fmt.Println(resp.Text())

}

func Test_Post(t *testing.T) {
	sess := requests.New()
	req := requests.NewRequest("POST", "http://httpbin.org/post", nil).SetRetry(3)
	req.SetParams(map[string]interface{}{
		"a": "b",
		"c": 3,
		"d": []int{1, 2, 3},
	})

	if err := req.SetBody(`{"body":"QWER"}`); err != nil {
		fmt.Println(err)
	}
	// req.SetHeader("Content-Type", "application/json")
	resp, err := sess.DoRequest(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resp.Text())
}
