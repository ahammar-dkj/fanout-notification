
# Fan-out notification via SNS and SQS

```
                        /--> SQS queue A --> consumer A (polling)
producer --> SNS topic <
                        \--> SQS queue B --> consumer B (polling)
```

## Localstack

[Localstack](https://github.com/localstack/localstack) can be used for local development without accessing SNS and SQS web services. 

Install Localstack:

```
pip install localstack awscli-local
```

Start SNS and SQS services with Localstack:

```
SERVICES=sns,sqs localstack start
```

Create SNS topic, SQS queues and setup subscription:

```
awslocal sns create-topic --name my_topic
awslocal sqs create-queue --queue-name my_queue_a
awslocal sqs create-queue --queue-name my_queue_b
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:my_topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:my_queue_a
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:my_topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:my_queue_b
```

## Producer

Run producer:

```
go run cmd/producer/main.go
```

## Consumer

Run first consumer:

````
go run cmd/consumer/main.go
````

Run second consumer:

```
go run cmd/consumer/main.go -queue-url http://localhost:4576/queue/my_queue_b
```
