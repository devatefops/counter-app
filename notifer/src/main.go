package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"
)

func main() {
	// Environment variables from docker-compose
	counterSvcHost := os.Getenv("COUNTER_SVC_HOST")
	smtpServer := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	emailTo := os.Getenv("EMAIL_TO")
	checkInterval := os.Getenv("CHECK_INTERVAL")

	intervalDuration, err := time.ParseDuration(checkInterval)
	if err != nil {
		log.Fatalf("incorrect interval: %v", err)
	}

	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("invalid SMTP port: %v", err)
	}

	for {
		checkCounterAndNotify(counterSvcHost, smtpServer, smtpPortInt, smtpUser, smtpPass, emailTo)
		log.Printf("Check complete. Sleeping for %s...", checkInterval)
		time.Sleep(intervalDuration)
	}
}

func checkCounterAndNotify(counterHost, smtpServer string, smtpPort int, smtpUser, smtpPass, emailTo string) {
	resp, err := http.Get("http://" + counterHost)
	if err != nil {
		log.Printf("Failed to get counter value: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Counter service returned non-OK status: %s", resp.Status)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	count, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		log.Printf("Failed to parse counter value: %v", err)
		return
	}

	log.Printf("Current counter value is %d", count)

	if count < 10 {
		log.Println("Counter is less than 10, skipping email.")
		return
	}

	body := fmt.Sprintf("The counter has reached a value of %d.", count)

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", smtpUser)
	message.SetHeader("To", emailTo)
	message.SetHeader("Subject", "Counter Alert!")

	// Set email body
	message.SetBody("text/html", `
        <html>
            <body>
               `+body+`
            </body>
        </html>
    `)
	// Set up the SMTP dialer
	dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPass)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	log.Println("Notification email sent successfully!")
}
