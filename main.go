package main

import (
	"fmt"

	"app/app/controllers"
	"app/app/models"
)

func main() {
	fmt.Println(models.Db)

	controllers.StartMainServer()

}
