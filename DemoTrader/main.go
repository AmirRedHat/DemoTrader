package main

import (
	// "local/models"
	// "os"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	userRouter := router.Group("/users")
	{
		//userRouter.GET("/:id", RetrieveUserView)
		userRouter.GET("/", RetrieveAllUsersView)
		userRouter.POST("/", CreateUserView)
		userRouter.POST("/login", LoginUserView)
		userRouter.POST("/token", TokenUserView)
	}

	tradeRouter := router.Group("/trade")
	{
		tradeRouter.GET("/user", RetriveTrade)
		tradeRouter.POST("/open", CreateTrade)
		tradeRouter.PUT("/close/:id", UpdateTrade)
	}

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
