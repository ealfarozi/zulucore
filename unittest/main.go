package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ealfarozi/zulucore/structs"
)

func main() {
	var usr = structs.User{}
	fmt.Println(usr)
	url := "http://localhost:8000/api/v1/login"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"username": "admin@superadmin.com", "password": "12345678"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
