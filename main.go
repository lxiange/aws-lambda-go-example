package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tb "gopkg.in/tucnak/telebot.v2"
	"encoding/json"
	"fmt"
	"os"
)

var TelegramToken = os.Getenv("telegram_token")

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("token:", TelegramToken)
	b, err := tb.NewBot(tb.Settings{
		Token: TelegramToken,
	})
	if err != nil {
		fmt.Println("create bot error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "error",
		}, fmt.Errorf("create bot error: %s\n", err)
	}

	update := tb.Update{}
	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		fmt.Printf("parse request error, request: %#v\n", request)
	}
	fmt.Printf("request update: %#v\n", update)
	if update.Message.FromGroup() {
		if update.Query != nil {
			b.Send(update.Message.Chat, "小君君真好看！")
		}
	} else {
		b.Send(update.Message.Chat, "小君君真好看！")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "ok",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
