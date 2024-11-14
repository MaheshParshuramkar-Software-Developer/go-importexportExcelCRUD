package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
)

// UploadedFileStatus : to check the uploaded file status
func UploadedFileStatus(c *gin.Context) {
	var uploadedFile models.UploadedFileStatus
	if err := c.ShouldBindJSON(&uploadedFile); err != nil {
		logger.Log.Debug("error valid json not provided ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to fetch data",
			"error":  "filename is required",
		})
		return
	}

	redisData, err := dbs.RClient.Get(uploadedFile.Filename)
	if err != nil {
		if err.Error() == "redis: nil" {
			// find from mysql table
			checkForStatusInDB(c, uploadedFile)
			return
		}
		logger.Log.Debug("error while getting data from redis ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to fetch data",
		})
		return
	}

	if redisData != "" {
		if redisData == "true" {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"result": "file was uploaded successfully",
				"error":  "",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"result": "file was not uploaded",
			"error":  "file validation failed",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "failure",
		"result": "Failed to fetch data",
		"error":  "Internal Server Error",
	})
}

func checkForStatusInDB(c *gin.Context, uploadedFile models.UploadedFileStatus) {
	result, err := dbs.GetUploadedFileStatus(uploadedFile.Filename)
	if err != nil {
		if err.Error() == "record not found" {
			logger.Log.Debug("no record found in db", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"result": "no records found for this filename",
				"error":  "",
			})
			return
		}
		logger.Log.Debug("error while getting data from db ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to fetch data",
			"error":  "error while getting data from database",
		})
		return
	}
	if result.Status == "true" {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"result": "file was uploaded successfully",
			"error":  "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": "file was not uploaded",
		"error":  "file validation failed",
	})
}

func uploadedFileStatusUpdate(filename string, isUploaded bool) {
	isUploadedStr := fmt.Sprintf("%t", isUploaded)
	updateFileStatKeyTll := time.Duration(Config.RedisConf.UpdateFileStatKeyTtl)
	if err := dbs.RClient.Set(filename, isUploadedStr, time.Hour*updateFileStatKeyTll); err != nil {
		logger.Log.Debug("error while storing data from redis ", err.Error())
	}
	uploadedFile := &models.UploadedFileStatus{
		Filename: filename,
		Status:   isUploadedStr,
	}

	if err := dbs.InsertFileStatus(uploadedFile); err != nil {
		logger.Log.Debug("error while storing data into db ", err.Error())
	}
}
