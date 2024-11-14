package controller

import (
	"encoding/json"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateData(c *gin.Context) {
	var updateDataUri models.UpdateDataUri
	if err := c.ShouldBindUri(&updateDataUri); err != nil {
		logger.Log.Debug("error id is not provided ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to update record",
			"error":  err.Error(),
		})
	}
	var updateData models.UpdateData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Log.Debug("invalid json body provided ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to update record",
			"error":  err.Error(),
		})
		return
	}

	employeeDataUp := &models.Employee{
		Id:        updateDataUri.Id,
		FirstName: updateData.FirstName,
		LastName:  updateData.LastName,
		Company:   updateData.Company,
		Address:   updateData.Address,
		City:      updateData.City,
		Country:   updateData.Country,
		Postal:    updateData.Postal,
		Phone:     updateData.Phone,
		Email:     updateData.Email,
		Web:       updateData.Web,
	}

	if err := dbs.UpdateEmployee(employeeDataUp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to update record",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status": "success",
		"Result": "Data updated successfully",
	})
	go updateDataInRedis(employeeDataUp)
}

// updateDataInRedis : update the record in redis if it exists or else add that record
func updateDataInRedis(employeeDataUp *models.Employee) {
	redisKey := Config.RedisConf.RedisUploadedDataKey
	redisData, err := dbs.RClient.Get(redisKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.Log.Debug("redis key-value does not exist")
		} else {
			logger.Log.Debug("Error while getting data from redis", err)
		}
	}
	var employees []*models.Employee
	if redisData != "" {
		if err := json.Unmarshal([]byte(redisData), &employees); err != nil {
			logger.Log.Debug("Error while unmarshalling data from redis", err)
		} else {
			var idFound bool
			for i := 0; i < len(employees); i++ {
				if employees[i].Id == employeeDataUp.Id {
					if employeeDataUp.FirstName != "" {
						employees[i].FirstName = employeeDataUp.FirstName
					}
					if employeeDataUp.LastName != "" {
						employees[i].LastName = employeeDataUp.LastName
					}
					if employeeDataUp.Company != "" {
						employees[i].Company = employeeDataUp.Company
					}
					if employeeDataUp.Address != "" {
						employees[i].Address = employeeDataUp.Address
					}
					if employeeDataUp.City != "" {
						employees[i].City = employeeDataUp.City
					}
					if employeeDataUp.Country != "" {
						employees[i].Country = employeeDataUp.Country
					}
					if employeeDataUp.Postal != "" {
						employees[i].Postal = employeeDataUp.Postal
					}
					if employeeDataUp.Phone != "" {
						employees[i].Phone = employeeDataUp.Phone
					}
					if employeeDataUp.Email != "" {
						employees[i].Email = employeeDataUp.Email
					}
					idFound = true
				}
			}
			if idFound {
				insertIntoRedis(employees)
				return
			}
		}
	}

	employee, err := dbs.GetEmployee(employeeDataUp.Id)
	if err != nil {
		logger.Log.Debug("Error while getting data from db", err)
	}
	employees = append(employees, &employee)
	insertIntoRedis(employees)
}
