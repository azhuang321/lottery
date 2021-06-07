package routes

import (
	"github.com/kataras/iris/v12/mvc"
	bootstrap "lottery/bootrap"
	"lottery/services"
	"lottery/web/controllers"
	"lottery/web/middleware"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := services.NewUserService()
	giftService := services.NewGiftService()
	codeService := services.NewCodeService()
	resultService := services.NewResultService()
	userdayService := services.NewUserdayService()
	blackipService := services.NewBlackipService()

	index := mvc.New(b.Party("/"))
	index.Register(
		userService,
		giftService,
		codeService,
		resultService,
		userdayService,
		blackipService,
	)
	index.Handle(new(controllers.IndexController))

	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(
		userService,
		giftService,
		codeService,
		resultService,
		userdayService,
		blackipService,
	)
	admin.Handle(new(controllers.AdminController))

	adminGift := mvc.New(b.Party("/gift"))
	adminGift.Router.Use(middleware.BasicAuth)
	adminGift.Register(
		giftService,
	)
	adminGift.Handle(new(controllers.AdminGiftController))

	adminCode := admin.Party("/code")
	adminCode.Register(codeService)
	adminCode.Handle(new(controllers.AdminCodeController))
}
