# Ngrok API Lambda Deployment

This repository provides the steps to build and package a Go application for deployment on AWS Lambda.

Blog Post: [Setting Up an AWS Lambda Triggered by API Gateway to Send Emails via SES](https://www.mindofguru.com/code/aws/2024/05/21/aws-lambda-api-ses.html)

## Build and Package

1. **Build the Go Application for Linux**:
   
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
   ```

2. **Create a Deployment Package**:
   
   ```bash
   zip lambda-handler-final.zip bootstrap
   ```

## Deployment

1. **Upload the `deployment.zip` file to AWS Lambda**:
    - In the AWS Management Console, go to the Lambda service.
    - Create a new Lambda function or update an existing one.
    - Choose "Upload from" and select ".zip file".
    - Upload the `deployment.zip` file.

2. **Set Environment Variables and Execution Role**:
    - Configure the required environment variables in your Lambda function.
    - Ensure your Lambda function’s execution role has the necessary permissions.
    
```