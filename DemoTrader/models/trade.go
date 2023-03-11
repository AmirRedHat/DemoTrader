package models

import (
	"fmt"
	"local/tools"
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type Trade struct {
	gorm.Model
	User 				User 		`json:"user_id" db:"user_id" gorm:"embedded"`
	OpenPrice 		float64 	`json:"open_price" db:"open_price"`
	ClosePrice		float64	`json:"close_price" db:"close_price"`
	Coin 				string 	`json:"coin" db:"coin"`
	Pair 				string 	`json:"pair" db:"pair"`
	PositionType	string 	`json:"position_type" db:"position_type"` // LONG or SHORT
	Quantity 		float64	`json:"quantity" db:"quantity"`
	TradeType 		string	`json:"trade_type" db:"trade_type"` // futures or spot
	Status 			string 	`json:"status" db:"status"`
	Profit 			float64	`json:"profit" db:"profit"`
}

func (trade *Trade) Migrate() {
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		fmt.Println("openening database error: ", err.Error())
		log.Fatal(err)
	}
	db.AutoMigrate(&Trade{})
	fmt.Println("trade model migrated !")
}

func (trade *Trade) Save() *Trade {

	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		fmt.Println("openening database error: ", err.Error())
		log.Fatal(err)
	}
	db.Create(&trade)
	return trade
}

func (trade *Trade) Read(id int) Trade {
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var targetTrade Trade
	db.Find(&targetTrade, []int{id})
	return targetTrade
}

func (trade *Trade) Update() {
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	id := trade.ID
	db.Model(Trade{}).Where("id = ?", id).Updates(trade)
}

func (trade *Trade) Validate() bool {
	fmt.Println("not implemented")
	if trade.TradeType == "futures" || trade.TradeType == "spot" {
		fmt.Println("validation complete")
	}
	return true
}

func (trader *Trade) GetLastPrice(coin string) float64 {
	url := fmt.Sprintf("http://82.115.25.213:4000/last_price?broker=kucoin&symbol=%s", coin)
	response := tools.Request(url, "GET", nil)
	return response["data"].(float64)
}

func (trader *Trade) FilterTrades(filterParams map[string]interface{}) []Trade {
	// read all records that contain filters
	db, err := gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var targetTrade []Trade
	db.Where(filterParams).Order("CreatedAt").Find(&targetTrade)
	return targetTrade
}

func (trader *Trade) ProfitCalculator() float64 {
	var profitPercent float64
	if trader.PositionType == "LONG" {
		profitPercent = ((trader.ClosePrice - trader.OpenPrice) / trader.Quantity) * 100
	} else if trader.PositionType == "SHORT" {
		profitPercent = ((trader.OpenPrice - trader.ClosePrice) / trader.Quantity) * 100
	} else {
		log.Fatal("invalid position type")
	}
	
	fmt.Println("[+] Profit calculated")
	return profitPercent
}