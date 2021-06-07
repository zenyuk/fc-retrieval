package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/utest"
)

func ExampleGetJSON_01() {
	var v01 int
	err := GetJSON("", &v01)
	fmt.Println(v01, err)
	// Output: 0 Get "": unsupported protocol scheme ""
}

func ExampleGetJSON_02() {
	v01 := make(map[string]string)

	httpClient = utest.NewTestClientString(`{"k1":"v1"}`)
	err := GetJSON("", &v01)
	httpClient = &http.Client{}

	fmt.Println(v01, err)
	// Output: map[k1:v1] <nil>
}

func ExampleSendJSON_01() {
	err := SendJSON("", 1)
	fmt.Println(err)
	// Output: Post "": unsupported protocol scheme ""
}

func ExampleSendJSON_02() {
	httpClient = utest.NewTestClientString(``)
	err := SendJSON("", 1)
	httpClient = &http.Client{}
	fmt.Println(err)
	// Output: <nil>
}

func ExampleSendMessage_01() {
	var req fcrmessages.FCRMessage
	var resp *fcrmessages.FCRMessage
	resp, err := SendMessage("", &req)
	fmt.Println(resp, err)
	// Output: <nil> Post "http:///v1": http: no Host in request URL
}

func ExampleSendMessage_02() {
	var req fcrmessages.FCRMessage
	var resp *fcrmessages.FCRMessage

	httpClient = utest.NewTestClientString(``)
	resp, err := SendMessage("", &req)
	httpClient = &http.Client{}

	fmt.Println(resp, err)
	// Output: &{0 0 [] [] } <nil>
}

func ExampleSendMessage_03() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()
	tsHostPort := ts.URL[7:]
	var req fcrmessages.FCRMessage
	var resp *fcrmessages.FCRMessage
	resp, err := SendMessage(tsHostPort, &req)
	fmt.Println(resp, err)
	// Output:
	// <nil> SendMessage receive error code: 404 Not Found
}

func ExampleSendMessage_04() {
	var req fcrmessages.FCRMessage
	var resp *fcrmessages.FCRMessage

	httpClient = utest.NewTestClientError()
	resp, err := SendMessage("", &req)
	httpClient = &http.Client{}

	fmt.Println(resp, err)
	// Output: &{0 0 [] [] } test error
}
