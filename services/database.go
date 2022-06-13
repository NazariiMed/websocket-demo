package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"time"
	"websocket-demo/internal/models"
)

type Database struct {
	config *models.Config
	pool   *pgxpool.Pool
}

func NewDatabase(config *models.Config) *Database {
	return &Database{
		config: config,
	}
}

func (this *Database) Start() {
	dbPool, err := pgxpool.Connect(context.Background(), this.config.DatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	this.pool = dbPool
}

func (this *Database) WriteTrade(trade models.Event) {
	//todo: at least 2 thing can be improved here:
	//1. save events as batch
	//2. replace plain sql with better approach like ORM or sqlx like.
	_, err := this.pool.Exec(context.Background(),
		"Insert into trades (evt, pair, tid, p, q, b, a, tt) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		time.UnixMilli(trade.EventTime), trade.Pair, trade.TradeId, trade.TradePrice,
		trade.TradeVolume, trade.ByOrderId, trade.SellOrderId, time.UnixMilli(trade.TradeCompletedTime))
	if err != nil {
		fmt.Printf("DB write error: ", err)
	}
}

//todo: add WriteTicker func similar as one for trades.
