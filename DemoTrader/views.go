package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"local/models"
	"net/http"
	"log"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func returnData(dataBytes io.Reader) map[string]interface{} {
	data := make(map[string]interface{})
	posted_data, err := ioutil.ReadAll(dataBytes)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(posted_data, &data)
	return data
}

func RetrieveAllUsersView(context *gin.Context) {
	var usr models.User
	usersList := usr.Read(0)
	context.AbortWithStatusJSON(200, usersList)
}

func CreateUserView(context *gin.Context) {
	var usr models.User

	err := context.BindJSON(&usr)
	if err != nil {
		fmt.Println(err)
		context.AbortWithStatus(400)
	} else {
		savedUser := usr.Save()
		context.IndentedJSON(http.StatusCreated, savedUser)
	}
}

func LoginUserView(context *gin.Context) {
	var usr models.User
	dataBytes := context.Request.Body
	data := returnData(dataBytes)
	email := data["email"].(string)
	password := data["password"].(string)
	isValid := usr.Validate(email, password)
	if isValid {
		context.AbortWithStatus(200)
	} else {
		context.AbortWithStatus(401)
	}
}

func TokenUserView(context *gin.Context) {
	context.JSON(http.StatusAccepted, "somestring")
}

func RetriveTrade(context *gin.Context) {
	var trd models.Trade
	lastPrice := trd.GetLastPrice("ETHUSDT")
	fmt.Println(lastPrice)
	context.AbortWithStatus(200)
}

func CreateTrade(context *gin.Context) {
	var trd models.Trade
	data := returnData(context.Request.Body)

	trd.User.UserId = int(data["user_id"].(float64))
	trd.Coin = data["coin"].(string)
	symbol := fmt.Sprintf("%s-%s", data["coin"].(string), data["pair"].(string))
	trd.OpenPrice = trd.GetLastPrice(symbol)
	trd.ClosePrice = 0
	trd.Pair = data["pair"].(string)
	trd.PositionType = data["position_type"].(string)
	trd.Quantity = data["quantity"].(float64)
	trd.TradeType = data["trade_type"].(string)
	trd.Status = "open"
	trd.Profit = 0
	trd.Save()

	context.AbortWithStatus(201)
}


func UpdateTrade(context *gin.Context) {
	var trd models.Trade
	id := context.Param("id")
	tradeId, err := strconv.Atoi(id)
	trd = trd.Read(tradeId)
	if err != nil {
		context.AbortWithStatus(400)
	}
	trd = trd.Read(tradeId)
	symbol := fmt.Sprintf("%s-%s", trd.Coin, trd.Pair)
	closePrice := trd.GetLastPrice(symbol)

	trd.ClosePrice = closePrice
	trd.Profit = trd.ProfitCalculator()

	trd.Update()

	context.AbortWithStatus(201)
}