package main

import (
	"github.com/roskeys/app/controller"
	"github.com/roskeys/app/utils"
)

func main() {
	utils.InitMainRouter()
	controller.InitAuthController()
	utils.MainRouter.Run(utils.SERVER_ADDRESS)
}
