package main

import (
	"SeedBot/core"
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	fmt.Println(`
  /$$$$$$                            /$$ /$$$$$$$              /$$    
 /$$__  $$                          | $$| $$__  $$            | $$    
| $$  \__/  /$$$$$$   /$$$$$$   /$$$$$$$| $$  \ $$  /$$$$$$  /$$$$$$  
|  $$$$$$  /$$__  $$ /$$__  $$ /$$__  $$| $$$$$$$  /$$__  $$|_  $$_/  
 \____  $$| $$$$$$$$| $$$$$$$$| $$  | $$| $$__  $$| $$  \ $$  | $$    
 /$$  \ $$| $$_____/| $$_____/| $$  | $$| $$  \ $$| $$  | $$  | $$ /$$
|  $$$$$$/|  $$$$$$$|  $$$$$$$|  $$$$$$$| $$$$$$$/|  $$$$$$/  |  $$$$/
 \______/  \_______/ \_______/ \_______/|_______/  \______/    \___/  
`)
	fmt.Println(`ρσωєяє∂ ву : нσℓу¢αη`)

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	core.ProcessBot(config.Default())
}
