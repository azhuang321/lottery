package main

import (
	"fmt"
	bootstrap "lottery/bootrap"
	"lottery/web/middleware/identity"
	"lottery/web/routes"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	//初始化应用
	app := bootstrap.New("go抽奖系统", "azhuang")
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)
	return app
}
func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
