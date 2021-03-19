package aws

import (
	"github.com/BrunoDM2943/go-todo-lambda/internal/cdi"
	"github.com/BrunoDM2943/go-todo-lambda/internal/handler/function"
	"github.com/aws/aws-lambda-go/lambda"
)

func StartLambda() {
	handler := function.NewLambdaHandler(cdi.GetTodoService())
	handler.BuildRoutes()
	lambda.Start(handler.HandleRequest)
}
