# ghe-repo-checker

Go program that collect all private microservices repositories in Github Enterprise, stores them in a DynamoDB database and sends notifications for every repo created or deleted. This program is intended to work in an AWS Lambda function, but you can run it locally with no problem
