package main

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
func createMessageWithAttachment(from string, to string, subject string, content string, fileDir string, fileName string) gmail.Message {

	var message gmail.Message

	// read file for attachment purpose
	// ported from https://developers.google.com/gmail/api/sendEmail.py

	fileBytes, err := ioutil.ReadFile(fileDir + fileName)
	if err != nil {
		log.Fatalf("Unable to read file for attachment: %v", err)
	}

	fileMIMEType := http.DetectContentType(fileBytes)

	// https://www.socketloop.com/tutorials/golang-encode-image-to-base64-example
	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32, "alphanum")

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + to + "\n" +
		"from: " + from + "\n" +
		"subject: " + subject + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		content + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		fileData +
		"--" + boundary + "--")

	// see https://godoc.org/google.golang.org/api/gmail/v1#Message on .Raw
	// use URLEncoding here !! StdEncoding will be rejected by Google API

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	return message
}
func sendMessage(service *gmail.Service, userID string, message gmail.Message) {
	_, err := service.Users.Messages.Send(userID, &message).Do()
	if err != nil {
		log.Fatalf("Unable to send message: %v", err)
	} else {
		log.Println("Email message sent!")
	}

}
func main() {
	bCredentialsAdmin, err := ioutil.ReadFile("credentials3.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(bCredentialsAdmin, gmail.GmailComposeScope)
	ctx := context.Background()
	ts := conf.TokenSource(ctx)
	log.Println(ts)
	conf.Subject = "pin2041to@pintest.page"

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := conf.Client(oauth2.NoContext)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	msgContent := `Hello!
	This is a test email send via Gmail API
	Good Bye!`
	messageWithAttachment := createMessageWithAttachment("phutthichod.t@ku.th", "pin2041to@gmail.com, pin2017create@outlook.com", "Email WITH ATTACHMENT from GMail API", msgContent, "./", "img.pdf")
	sendMessage(srv, "me", messageWithAttachment)
}
