#!/usr/bin/env bash -x

awslocal sns create-topic --name my_topic
awslocal sqs create-queue --queue-name my_queue_a
awslocal sqs create-queue --queue-name my_queue_b
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:my_topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:my_queue_a
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:my_topic --protocol sqs --notification-endpoint arn:aws:sns:us-east-1:000000000000:my_queue_b
