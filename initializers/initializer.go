package initializers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-importexportExcelCRUD/config"
	"github.com/go-importexportExcelCRUD/controller"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
)

var server *http.Server

var Config *models.Configurations

func LoadConfig() *models.Configurations {
	Config = config.GetConfig()
	logger.InitializeLogger(Config.LogLevel, Config.LogPath)
	controller.LoadConfigForCtl(Config)
	return Config
}

func StartServer(r *gin.Engine) {
	server = &http.Server{
		Addr:         Config.ServerConf.Host + ":" + Config.ServerConf.Port,
		ReadTimeout:  time.Duration(Config.ServerConf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.ServerConf.WriteTimeout) * time.Second,
		Handler:      r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("error while listen :", err)
		}
	}()
}

func InitializeConnections() {
	dbs.ConnectDb(&Config.MySqlConf)
	dbs.SyncDb()
	dbs.NewRedisClient(&Config.RedisConf)
	logger.Log.Info("Connections made Successfully")
}

// StopServices : stop signal received closing the opened connections
func StopServices() {
	dbs.RClient.Close()
	stopServer()
}

func stopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Error while shutting down server", err.Error())
	}
	logger.Log.Debug("Server gracefully stopped")
}
