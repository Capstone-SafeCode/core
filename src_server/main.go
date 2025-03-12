package main

import (
	"test_capstone/src_server/controllers"
	"test_capstone/src_server/database"
	"test_capstone/src_server/routes"
)

func main() {
	database.ConnectDB()

	controllers.InitCollections()

	router := routes.SetupRouter()
	router.Run(":8080")
}
