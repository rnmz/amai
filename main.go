package main

import app "amai/blog/app/config"

func main() {
	loadEnvErr := app.InitEnvVariables()
	if loadEnvErr != nil {
		panic(loadEnvErr)
	}

	checkEnvErr := app.CheckEnvParams()
	if checkEnvErr != nil {
		panic(checkEnvErr)
	}

	app.ShowStartMessage()
	app.StartServer()
}
