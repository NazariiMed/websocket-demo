package services

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
	"websocket-demo/internal/models"
)

type Websocket struct {
	config   *models.Config
	database *Database
	conn     *websocket.Conn
}

func NewWebsocket(config *models.Config, database *Database) *Websocket {
	return &Websocket{
		config:   config,
		database: database,
	}
}

func (this *Websocket) Start(wg *sync.WaitGroup) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", this.config.WsUrl)

	c, _, err := websocket.DefaultDialer.Dial(this.config.WsUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	this.conn = c
	defer c.Close()
	wg.Done()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			//todo: handle errors properly instead of just logging it. Restart all subscriptions.
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			var trade models.Event
			err = json.Unmarshal(message, &trade)
			if err != nil {
				log.Printf("Err:", err)
			}
			//todo handle ticker events e.g.
			//if trade.Event == "ticker" {
			//	this.database.WriteTrade(trade)
			//}
			if trade.Event == "trade" {
				this.database.WriteTrade(trade)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (this *Websocket) Subscribe(pair string) {
	err := this.conn.WriteMessage(websocket.TextMessage, []byte("{\"method\":\"SUBSCRIBE\",\"params\":[\""+pair+"\"],\"id\":1}"))
	if err != nil {
		log.Println("error write message:", err)
	}
}
