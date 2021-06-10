package utest

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ExampleGetFreePort_01() {
	port := GetFreePort()
	fmt.Println(port != "")
	// Output: true
}

func ExampleNewTestClientString() {
	c := NewTestClientString("OK")
	fmt.Printf("%T\n", c)

	r, err := c.Get("")
	fmt.Println(r.Body)
	fmt.Println(err)
	// Output:
	// *http.Client
	// {OK}
	// <nil>
}

func ExampleNewTestClientError() {
	c := NewTestClientError()
	fmt.Printf("%T\n", c)

	r, err := c.Get("")
	fmt.Println(err)

	bytes, err := ioutil.ReadAll(r.Body)
	fmt.Println(bytes)
	fmt.Println(err)
	// Output:
	// *http.Client
	// <nil>
	// []
	// test error
}

func ExampleNewTestClient() {
	fn := func(req *http.Request) *http.Response {
		return &http.Response{}
	}
	c := NewTestClient(fn)
	fmt.Printf("%T\n", c)
	fmt.Printf("%T\n", c.Transport)
	// Output:
	// *http.Client
	// utest.RoundTripFunc
}
