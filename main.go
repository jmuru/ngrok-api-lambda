package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)


type EndpointList struct {
	Endpoints   []Endpoint `json:"endpoints"`
	URI         string     `json:"uri"`
	NextPageURI string     `json:"next_page_uri"`
}

type Endpoint struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PublicURL string    `json:"public_url"`
	Proto     string    `json:"proto"`
	HostPort  string    `json:"hostport"`
	Type      string    `json:"type"`
	Tunnel    Tunnel    `json:"tunnel"`
}

type Tunnel struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

func handler(ctx context.Context) error {
	apiKey := os.Getenv("API_KEY")
	// Create an HTTP client
	client := &http.Client{}

	// Set up the GET request with the Authorization and Ngrok-Version headers
	req, err := http.NewRequest("GET", "https://api.ngrok.com/endpoints", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Ngrok-Version", "2")

	// Make the request and handle the response
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var endpointList EndpointList
	err = json.NewDecoder(resp.Body).Decode(&endpointList)
	if err != nil {
		return err
	}

	// Collect public URLs
	publicURLs := ""
	for _, endpoint := range endpointList.Endpoints {
		publicURLs += fmt.Sprintf("Public URL: %s\n", endpoint.PublicURL)
	}

	return sendEmail(publicURLs)
}

func sendEmail(publicURLs string) error {
	// Get the SES region and the sender and recipient emails from environment variables
	region := os.Getenv("REGION")
	sender := os.Getenv("SENDER")
	recipient := os.Getenv("RECIPIENT")
	if region == "" || sender == "" || recipient == "" {
		return fmt.Errorf("AWS_REGION, SENDER_EMAIL, or RECIPIENT_EMAIL environment variable is not set")
	}

	// Create a new SES session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return fmt.Errorf("Error creating AWS session: %v", err)
	}

	svc := ses.New(sess)

	// Set the email parameters
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(publicURLs),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Ngrok Public URLs"),
			},
		},
		Source: aws.String(sender),
	}

	// Send the email
	_, err = svc.SendEmail(input)
	if err != nil {
		return fmt.Errorf("Error sending email: %v", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
