package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-importexportExcelCRUD/controller"
	"github.com/go-importexportExcelCRUD/models"
)

func SetupRouter(cfg *models.Configurations) *gin.Engine {
	r := gin.Default()
	r.GET(cfg.ServerConf.Prefix+"/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST(cfg.ServerConf.Prefix+"/uploadData", controller.UploadFile)
	r.GET(cfg.ServerConf.Prefix+"/viewData", controller.ViewData)
	r.PUT(cfg.ServerConf.Prefix+"/updateData/:id", controller.UpdateData)
	r.DELETE(cfg.ServerConf.Prefix+"/deleteData/:id", controller.DeleteData)

	r.POST(cfg.ServerConf.Prefix+"/fileUploadedStatus", controller.UploadedFileStatus)
	return r
}
