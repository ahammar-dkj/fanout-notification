package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsQueueURL string

func init() {
	const (
		defaultQueueURL = "http://localhost:4576/queue/my_queue_a"
	)
	flag.StringVar(&sqsQueueURL, "queue-url", defaultQueueURL, "SQS queue URL")
	flag.Parse()
}

func main() {
	fmt.Println("Queue URL:", sqsQueueURL)

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4576"),
		Region:   aws.String("us-east-1"),
	}))
	svc := sqs.New(sess)

	ticker := time.NewTicker(10 * time.Second)
	for t := range ticker.C {
		fmt.Println("Polling at", t)
		receiveMessages(svc)
	}
}

func receiveMessages(svc *sqs.SQS) {
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &sqsQueueURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(20),
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		log.Println("Failed to receive messages", err)
	}
	defer deleteMessages(svc, result.Messages)

	fmt.Printf("Received %d messages\n", len(result.Messages))
	for _, m := range result.Messages {
		fmt.Println(m)
	}
}

func deleteMessages(svc *sqs.SQS, messages []*sqs.Message) {
	if len(messages) == 0 {
		return
	}

	var entries []*sqs.DeleteMessageBatchRequestEntry
	for _, m := range messages {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            m.MessageId,
			ReceiptHandle: m.ReceiptHandle,
		})
	}
	result, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: &sqsQueueURL,
	})
	if err != nil {
		log.Println("Failed to delete message batch", err)
	}

	fmt.Printf("Successfully deleted %d messages\n", len(result.Successful))
}
