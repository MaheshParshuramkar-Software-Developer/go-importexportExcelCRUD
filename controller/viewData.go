package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
)

func ViewData(c *gin.Context) {
	redisKey := Config.RedisConf.RedisUploadedDataKey
	redisData, err := dbs.RClient.Get(redisKey)
	if err != nil {
		// data not found in redis fetching data from mysql
		if err.Error() == "redis: nil" {
			employees, err := dbs.GetAllEmployee()
			if err != nil {
				logger.Log.Debug("error while getting data from from database", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failure",
					"result": "failed to get data",
					"error":  "Internal Server Error!",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"result": employees,
				"error":  "",
			})
			// setting the same data in redis cache also for future use
			go storeDataInCache(employees)
			return
		}
		logger.Log.Debug("error while getting data from from redis", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"result": "failed to get data",
			"error":  "Internal Server Error!",
		})
		return
	}

	// data found in redis
	if redisData != "" {
		var viewData []models.ViewData
		err := json.Unmarshal([]byte(redisData), &viewData)
		if err != nil {
			logger.Log.Debug("error while unmarshalling data", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failure",
				"result": "failed to get data",
				"error":  "Internal Server Error!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"result": viewData,
			"error":  "",
		})
		return
	}
	// if code blocks comes here in any particular case
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "failure",
		"result": "failed to get data",
		"error":  "Internal Server Error!",
	})
}

func storeDataInCache(employees []models.ViewData) {
	redisKey := Config.RedisConf.RedisUploadedDataKey
	redisKeyTtl := time.Duration(Config.RedisConf.RedisUploadedDataKeyTtl)
	redisData, err := json.Marshal(employees)
	if err != nil {
		fmt.Println("error while marshalling data", err.Error())
	}
	if err := dbs.RClient.Set(redisKey, string(redisData), time.Minute*redisKeyTtl); err != nil {
		logger.Log.Debug("error while setting data into redis", err.Error())
	}
}
