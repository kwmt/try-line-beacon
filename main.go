package main

import (
	"log"
	"net/http"
	"os"
	"io/ioutil"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req)
	})

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		//log.Println("/callback", req)
		body, _ := ioutil.ReadAll(req.Body)
		log.Println("body", string(body))


		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			log.Println("type", event.Type)
			if event.Type == linebot.EventTypeBeacon {
				//log.Println("type", event.Type)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Println("message", message)
					log.Println("linebot.NewTextMessage(message.Text)", linebot.NewTextMessage(message.Text))
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}