package mqtt

import (
	"bytes"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"io/ioutil"
	"net/http"
)

func GetConnections() {
	auth()

	res, err := http.Get(base.DefaultConfig.MqttServerHttp + "/api/v3/connections/")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func auth() {
	user := make(map[string]interface{})
	user["username"] = "admin"
	user["password"] = "public"

	tmp := `{"username":"admin", "password": "public"}`
	req := bytes.NewBuffer([]byte(tmp))

	res, err := http.Post(base.DefaultConfig.MqttServerHttp + "/api/v3/auth", "application/json",req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}