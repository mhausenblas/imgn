package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type uploadres struct {
	Status string `json:"status"`
}

func upload() (*uploadres, error) {
	ur := &uploadres{
		Status: "all good",
	}

	return ur, nil
}

func main() {
	lambda.Start(upload)
}
