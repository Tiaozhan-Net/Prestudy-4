package main

import (
	"fmt"
	"fe_exam/config"
	"fe_exam/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	srv := router.NewServer()

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("fail to init server: %s\n", err.Error())
		panic(err)
	}
}
