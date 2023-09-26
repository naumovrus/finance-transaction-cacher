package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/naumovrus/finance-transaction-api/internal/cache"
	"github.com/naumovrus/finance-transaction-api/internal/service"
)

type Handler struct {
	services   *service.Service
	redisCache cache.TransactionCache
}

func NewHandler(services *service.Service, redisCache cache.TransactionCache) *Handler {
	return &Handler{
		services:   services,
		redisCache: redisCache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		finance := api.Group("/finance")
		{
			finance.POST("/create-wallet", h.createWallet)
			// finance.GET("/wallet/:id", h.getDataUser)
			finance.POST("/top-up", h.topUp)
			finance.POST("/take-out", h.takeOut)
			finance.POST("/send", h.send)

		}

	}
	return router
}
