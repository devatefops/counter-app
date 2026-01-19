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

	// Send a welcome email on startup
	sendWelcomeEmail(cfg)

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

	if count != 10 {
		log.Printf("Counter is not 10 (value: %d), skipping email.", count)
		return
	}

	// Send notification email
	subject := "Counter Alert!"
	body := fmt.Sprintf("The counter has reached the target value of %d.", count)
	if err := sendEmail(cfg, subject, body); err != nil {
		log.Printf("Failed to send notification email: %v", err)
	} else {
		log.Println("Notification email sent successfully!")
	}
}

func sendWelcomeEmail(cfg Config) {
	log.Println("Sending welcome email...")
	subject := "Notifier Service Started"
	body := "Welcome! The notifier service is running and will alert you when the counter reaches 10."
	if err := sendEmail(cfg, subject, body); err != nil {
		// Log the error but don't stop the service
		log.Printf("Failed to send welcome email: %v", err)
	} else {
		log.Println("Welcome email sent successfully!")
	}
}

// sendEmail is a helper function to send emails.
func sendEmail(cfg Config, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", cfg.SMTPUser) // The "From" address can be anything, but SMTP user is a good default
	message.SetHeader("To", cfg.EmailTo)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", `
        <html>
            <body>
               <p>`+body+`</p>
            </body>
        </html>
    `)
	dialer := gomail.NewDialer(cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	return dialer.DialAndSend(message)
}
