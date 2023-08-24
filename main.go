package main

import (
	"assignment-project/routers"
)

	var PORT = ":9090"


func main() {
	routers.StartServer().Run(PORT)
}