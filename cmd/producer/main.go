package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var snsTopicArn string

func init() {
	const (
		defaultTopic = "arn:aws:sns:us-east-1:000000000000:my_topic"
	)
	flag.StringVar(&snsTopicArn, "topic-arn", defaultTopic, "SNS topic ARN")
	flag.Parse()
}

func main() {
	fmt.Println("Topic ARN:", snsTopicArn)

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:4575"),
		Region:   aws.String("us-east-1"),
	}))
	svc := sns.New(sess)

	fmt.Printf("Press ENTER to send message")
	for {
		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		sendMessage(svc)
	}
}

func sendMessage(svc *sns.SNS) {
	data := struct {
		Message string `json:"msg"`
	}{
		fmt.Sprintf("Current time is: %v\n", time.Now()),
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Print("Failed to marshal JSON", err)
	}

	result, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(bytes)),
		TopicArn: &snsTopicArn,
	})
	if err != nil {
		log.Print("Failed to send message", err)
	}

	fmt.Print("Message sent: ", result.MessageId)
}
