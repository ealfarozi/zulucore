package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ealfarozi/zulucore/structs"
)

func main() {
	url := "http://localhost:8000/api/v1/login"
	fmt.Println("URL:>", url)

	var usr = structs.User{Username: "admin@superadmin.com", Password: "12345678"}
	fmt.Println("struct value: ", usr)

	requestByte, _ := json.Marshal(usr)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestByte))
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
