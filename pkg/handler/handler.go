package handler

import (
	"github.com/fancurson/toDoList/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		list := api.Group("/lists")
		{
			list.POST("/", h.createList)
			list.GET("/", h.getAllLists)
			list.GET("/:id", h.getListsById)
			list.PUT("/:id", h.updateLists)
			list.DELETE("/:id", h.deleteList)

			items := list.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("item")
		{
			items.GET("/:id", h.getItemsById)
			items.PUT("/:id", h.updateItems)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
