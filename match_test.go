package sinister

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func get1(hc *HC) {
	hc.MIME(ApplicationJSON)
	hc.Log("get1", WARN)
	hc.JSONS(http.StatusOK, "get1")
}
func post1(hc *HC) {
	hc.MIME(ApplicationJSON)
	hc.Log("post1", WARN)
	paramID := hc.Param("cacat")
	fmt.Println(paramID)
	hc.JSONS(http.StatusOK, "post1")
}

func TestSinister(t *testing.T) {
	sinister1 := New()
	sinister1.GET("/get", get1)
	sinister1.POST("/", get1)
	sinister1.POST("/post/[id]", post1)
	// sinister1.Start()

	get1 := httptest.NewServer(sinister1.router)
	defer get1.Close()
	res, err := http.Get(get1.URL + "/get")
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	fmt.Printf("%s", respBody)

	post1 := httptest.NewServer(sinister1.router)
	defer get1.Close()
	res2, err := http.Post(post1.URL+"/post/5", "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	respBody2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		log.Fatal(err)
	}
	res2.Body.Close()
	fmt.Printf("%s", respBody2)
}
func TestValidatepath(t *testing.T) {
	// formatted, params := validatePath("/v1/user/[id125125]", "GET")
	// wtfman("/v1/user/[ds]/[awdawd]/okkknigga/[s21]")
	a, b := validatePath("/v1/users/[id]", "GET")
	fmt.Println(a)
	fmt.Println(b)
	// t.Logf("formatted: %v params: %v\n", formatted, params)
	// t.Logf("%v", )

}
func BenchmarkValidate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// plm := wtfman("/v1/user/[ds]/[awdawd]/okkknigga/[s21]")
		// _ = plm
		a, b := validatePath("/v1/users/[id]", "GET")
		fmt.Println(a)
		fmt.Println(b)
	}
}

func BenchmarkRouter(b *testing.B) {
	b.ReportAllocs()

}

var pula = map[string]string{
	"/home": "GET",
	"/lol/[id]/[whaaat]/lol/okokokok/okokok/[aaa]": "POST",
	"/what": "DELETE",
}

func BenchmarkMatch(b *testing.B) {

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for k, v := range pula {
			a, b := validatePath(k, v)
			_ = a
			_ = b
		}
	}
}
