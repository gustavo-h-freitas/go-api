package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gustavo-h-freitas/web-service-gin/internals/database"
	"github.com/gustavo-h-freitas/web-service-gin/internals/routes"
)

func main() {
	router := gin.Default()
	db := database.GetClientDb()
	routes.DefineRoutes(router, db)
}
