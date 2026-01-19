package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"
)

// Config holds all configuration for the application.
type Config struct {
	CounterSvcHost   string
	SMTPServer       string
	SMTPPort         int
	SMTPUser         string
	SMTPPass         string
	EmailTo          string
	CheckInterval    time.Duration
}

func main() {
	smtpPortInt, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("invalid SMTP port: %v", err)
	}

	intervalDuration, err := time.ParseDuration(os.Getenv("CHECK_INTERVAL"))
	if err != nil {
		log.Fatalf("incorrect interval: %v", err)
	}

	cfg := Config{
		CounterSvcHost: os.Getenv("COUNTER_SVC_HOST"),
		SMTPServer:     os.Getenv("SMTP_HOST"),
		SMTPPort:       smtpPortInt,
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPass:       os.Getenv("SMTP_PASS"),
		EmailTo:        os.Getenv("EMAIL_TO"),
		CheckInterval:  intervalDuration,
	}

	for {
		checkCounterAndNotify(cfg)
		log.Printf("Check complete. Sleeping for %s...", cfg.CheckInterval)
		time.Sleep(cfg.CheckInterval)
	}
}

func checkCounterAndNotify(cfg Config) {
	resp, err := http.Get("http://" + cfg.CounterSvcHost)
	if err != nil {
		log.Printf("Failed to get counter value: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Counter service returned non-OK status: %s", resp.Status)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
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
	message.SetHeader("From", cfg.SMTPUser)
	message.SetHeader("To", cfg.EmailTo)
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
	dialer := gomail.NewDialer(cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	log.Println("Notification email sent successfully!")
}
