package mqtt

import (
	"encoding/base64"
	"fmt"
	"github.com/robertzml/Glaucus/base"
	"io/ioutil"
	"net/http"
)

func GetConnections() {
	// auth()

	client := &http.Client{}
	req, err := http.NewRequest("GET", base.DefaultConfig.MqttServerHttp + "/api/v3/connections/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth := "admin:public"
	authstr :=  base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", "Basic " + authstr)

	res, err := client.Do(req)
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
