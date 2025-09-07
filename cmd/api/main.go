package main

import (
	"kami/internal/boot"
	"kami/internal/config"
	"kami/internal/stmp"
)

func main() {
	config.LoadConfig()
	stmp.ConnectIMAP()
	go stmp.GetMailinbox()
	boot.Boot()
}
