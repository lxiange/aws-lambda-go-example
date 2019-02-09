package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tb "gopkg.in/tucnak/telebot.v2"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var TelegramToken = os.Getenv("telegram_token")

var bot, _ = tb.NewBot(tb.Settings{
	Token: TelegramToken,
})

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	update := tb.Update{}
	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		fmt.Printf("parse request error, request: %#v\n", request)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "error",
		}, fmt.Errorf("parse request error, request: %#v\n", request)
	}

	message := update.Message
	if message == nil {
		fmt.Println("Message is nil, skip")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "ok",
		}, nil
	}

	fmt.Printf("request message: %#v\n", message)

	fmt.Println(message.Entities, len(message.Entities))
	if message.FromGroup() {
		if strings.Contains(message.Text, "小君") {
			fmt.Println("hit!!!!!!!!!")
			bot.Send(update.Message.Chat, "小君君真好看！")
		}
	} else {
		bot.Send(update.Message.Chat, "小君君真好看！")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "ok",
	}, nil
}

func main() {
	fmt.Println("token:", TelegramToken)
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
