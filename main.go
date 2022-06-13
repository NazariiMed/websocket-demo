package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"websocket-demo/internal/models"
	"websocket-demo/internal/utils"
	"websocket-demo/services"
)

func main() {
	logrus.Info("Starting svc-api-gateway service...")

	cfg, err := models.ParseConfig()
	utils.PanicOnError(err)

	logrus.Infof("Configuration Loaded:\n\n %+v", cfg)

	database := services.NewDatabase(&cfg)
	database.Start()

	websocket := services.NewWebsocket(&cfg, database)
	var wg sync.WaitGroup
	wg.Add(1)
	go websocket.Start(&wg)
	wg.Wait()

	uniquePairs := getTradePairs(cfg.HttpBaseUrl)

	for _, e := range uniquePairs {
		//todo subscribe to ticker events
		//websocket.Subscribe(strings.ToLower(e) + "@ticker")
		websocket.Subscribe(strings.ToLower(e) + "@trade")
	}

	server := services.NewServer(&cfg)
	server.Start()
}

//todo: this method always return 1 value BTCUSDT which does not make sense to me.
//If there is another source of pairs in the API, then this method should be rewritten
//Also, it makes sense to separate all API calls to binance in the separate class
//Also, if pairs changes often we need to run it on schedule and add new subscriptions.
func getTradePairs(baseUrl string) []string {
	resp, err := http.Get(baseUrl + "vapi/v1/optionInfo")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var tradingPairs models.TradingPairsDataDTO
	err = json.Unmarshal(body, &tradingPairs)
	if err != nil {
		log.Printf("Error parsing trading pairs info: ", err)
	}

	keys := make(map[string]bool)
	uniquePairs := []string{}
	for _, entry := range tradingPairs.Data {
		if _, value := keys[entry.Underlying]; !value {
			keys[entry.Underlying] = true
			uniquePairs = append(uniquePairs, entry.Underlying)
		}
	}
	return uniquePairs
}
