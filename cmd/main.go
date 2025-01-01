package main

import (
	"fmt"
	"gestia/cmd/app"
	"log"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("gestia application successfully running!")

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	app, err := app.NewApp(logger)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.ListenAndServe("127.0.0.1:9090"))
}
