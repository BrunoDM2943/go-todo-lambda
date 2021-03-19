package aws

import (
	"github.com/BrunoDM2943/go-todo-lambda/internal/cdi"
	"github.com/BrunoDM2943/go-todo-lambda/internal/handler/function"
	"github.com/aws/aws-lambda-go/lambda"
)

func StartLambda() {
	lambda.Start(function.NewLambdaHandler(cdi.GetTodoService()).HandleRequest)
}
