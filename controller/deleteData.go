package controller

import (
	"encoding/json"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteData(c *gin.Context) {
	var deleteDataUri models.UpdateDataUri
	if c.ShouldBindUri(&deleteDataUri) != nil {
		logger.Log.Debug("error id is not provided ")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to delete data",
			"error":  "id is required",
		})
		return
	}

	if err := dbs.DeleteEmployee(deleteDataUri.Id); err != nil {
		logger.Log.Debug("error while deleting record from db ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to delete data",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": "record deleted successfully",
		"error":  "",
	})

	// delete from redis also
	go deleteDataFromRedis(deleteDataUri.Id)
}

func deleteDataFromRedis(id uint) {
	redisKey := Config.RedisConf.RedisUploadedDataKey
	redisData, err := dbs.RClient.Get(redisKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.Log.Debug("no records found in redis not deleteing from redis")
			return
		}
		logger.Log.Debug("error while getting data from redis", err.Error())
		return
	}

	if redisData != "" {
		var employees []models.ViewData
		if err := json.Unmarshal([]byte(redisData), &employees); err != nil {
			logger.Log.Debug("Error while unmarshalling data from redis", err)
		} else {
			var idFound bool
			for i := 0; i < len(employees); i++ {
				if employees[i].Id == id {
					employeesNew := append(employees[:i], employees[i+1:]...)
					employees = employeesNew
					idFound = true
					break
				}
			}

			if idFound {
				redisDataToSet, err := json.Marshal(employees)
				if err != nil {
					logger.Log.Debug("error while marshalling  the redis data ", err.Error())
				}
				redisKeyTtl := time.Duration(Config.RedisConf.RedisUploadedDataKeyTtl)
				// setting the updated data in redis if it exists
				if err := dbs.RClient.Set(redisKey, string(redisDataToSet), time.Minute*redisKeyTtl); err != nil {
					logger.Log.Debug("error while inserting the records into redis", err.Error())
				}
				return
			}
			logger.Log.Debug("no data found in redis to set for this id")
		}
	}
}
