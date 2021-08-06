package ip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func BoundInternetIP() (ipList [1]string, err error) {
	log.Println("This is a network request, please use it according to the actual situation.")
	responseClient, err := http.Get("http://ip.dhcp.cn/?ip")
	if err != nil {
		log.Println(err)
		return
	}
	defer responseClient.Body.Close()
	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	ipList[0] = clientIP
	return
}
