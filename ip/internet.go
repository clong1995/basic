package ip

import (
	"fmt"
	"io"
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
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(responseClient.Body)
	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	ipList[0] = clientIP
	return
}
