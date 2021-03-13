package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vandario/govms-ipg/app/http/controllers"
)

func ApiRoutes() *gin.Engine {
	apiRoutes := gin.Default()

	apiRoutes.POST("api/ipg/send", controllers.GatewayIpgSendHandler) //TODO: White house.
	apiRoutes.GET("ipg/:token", controllers.GatewayIpgRequestHandler) //TODO: White house.

	return apiRoutes
}
