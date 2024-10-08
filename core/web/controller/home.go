package controller

import (
	"net/http"

	"varconf/core/moudle/router"
	"varconf/core/service"
	"varconf/core/web/common"
)

type HomeController struct {
	common.Controller

	homeService *service.HomeService
}

func InitHomeController(s *router.Router, homeService *service.HomeService) *HomeController {
	homeController := HomeController{homeService: homeService}

	s.Get("/home/overall", homeController.overall)

	return &homeController
}

// GET /home/overall
func (_self *HomeController) overall(w http.ResponseWriter, r *http.Request, c *router.Context) {
	common.WriteSucceedResponse(w, _self.homeService.Overall())
}
