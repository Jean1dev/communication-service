package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DispatchSMS(to string, text string) {
	url := "https://vj44dv.api.infobip.com/sms/2/text/advanced"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"messages":[{"destinations":[{"to":"%s"}],"from":"ServiceSMS","text":"%s"}]}`, to, text))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", os.Getenv("INFOBIP_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
