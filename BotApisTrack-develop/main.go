package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"

	"github.com/aws/aws-lambda-go/lambda"
)

type Api struct {
	Url        string
	Xkey       string
	Key        string
	AnswerText string
	StatusCode int
	HttpMethod string
}

type Apis struct {
	List []Api
}

func testApis(scanApi *Api) bool {

	if scanApi.Url == "" {

		return false

	}

	client := &http.Client{

		CheckRedirect: func(req *http.Request, via []*http.Request) error {

			return http.ErrUseLastResponse

		},
	}

	req, _ := http.NewRequest("GET", scanApi.Url, nil)

	if scanApi.HttpMethod == "POST" {

		req, _ = http.NewRequest("POST", scanApi.Url, nil)

	}

	req.Header.Add("x-api-key", scanApi.Key)

	resp, err := client.Do(req)

	responseBody, _ := ioutil.ReadAll(resp.Body)

	if err == nil {
		if resp.StatusCode == scanApi.StatusCode || string(responseBody) == scanApi.AnswerText {

			fmt.Println(scanApi.Url + ": " + "ok")

		} else {

			fmt.Println("Lambda ", scanApi.Url, " returned error code: ", resp.StatusCode)

			SendMail(scanApi.Url, resp.StatusCode)

		}
	}
	defer resp.Body.Close()

	return true
}

func CheckLoop() {

	scanApi := []Api{
		{
			Url:        "https://97dd9d6yk4.execute-api.us-east-1.amazonaws.com/prod/",
			Key:        os.Getenv("KEY_SCRAPER"),
			AnswerText: "{\"message\":\"Empty URL\",\"status\":\"error\"}",
			HttpMethod: "GET",
			StatusCode: 400,
		},
		{
			Url:        "https://qw9dt71jg1.execute-api.us-east-1.amazonaws.com/mrshort/gxGXGs6",
			Key:        "",
			AnswerText: "<a href=\"https://google.com\">Found</a>.",
			HttpMethod: "GET",
			StatusCode: 302,
		},
		{
			Url:        "https://uax20edb40.execute-api.us-east-1.amazonaws.com/mrshort/",
			Key:        os.Getenv("KEY_MRSHORT"),
			AnswerText: "{\"Error\":\"Wrong URL format!\",\"URL\":\"\",\"status\":\"error\"}",
			HttpMethod: "GET",
			StatusCode: 400,
		},
		{
			Url:        "https://6dhm6gkofk.execute-api.us-east-1.amazonaws.com/qrcode/text",
			Key:        os.Getenv("KEY_QRCODE"),
			AnswerText: "{\"message\":\"Text is empty\",\"status\":\"error\"}",
			HttpMethod: "POST",
			StatusCode: 400,
		},
		{
			Url:        "https://u9406d69n8.execute-api.us-east-1.amazonaws.com/search/",
			Key:        os.Getenv("KEY_GOOGLESEARCH"),
			AnswerText: "{\"message\":\"Empty query\",\"status\":\"error\"}",
			HttpMethod: "GET",
			StatusCode: 400,
		},
	}

	a := Apis{scanApi}

	for _, scanApi := range a.List {

		testApis(&scanApi)

	}
}

func SendMail(url string, StatusCode int) {

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
	    fmt.Println("SMTP_PASSWORD environment variable is not set.")
	    return
	}
	
	auth := smtp.PlainAuth(
	    "",
	    "pushkin85.mil@gmail.com",
	    smtpPassword,
	    "smtp.gmail.com",
	)

	msg := fmt.Sprintf("To: pushkin85.mil@gmail.com\r\n"+
		"Subject: Lambda is down!\r\n"+
		"\r\n"+"Site %s returned error code %v", url, StatusCode)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"pushkin85.mil@gmail.com",
		[]string{"pushkin85.mil@gmail.com"},
		[]byte(msg),
	)

	if err != nil {

		fmt.Println(err)

	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// CheckLoop() // local start

	lambda.Start(CheckLoop) 

}
