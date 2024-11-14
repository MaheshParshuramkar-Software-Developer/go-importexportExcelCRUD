package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-importexportExcelCRUD/dbs"
	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
)

var Config *models.Configurations

func LoadConfigForCtl(config *models.Configurations) {
	Config = config
}

func UploadFile(c *gin.Context) {
	// single file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		logger.Log.Debug("error while reading the file ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to upload data",
			"error":  "invalid file !",
		})
		return
	}
	defer file.Close()

	// check for valid Excel file
	if !isValidExcelFile(fileHeader.Filename) {
		logger.Log.Debug("valid file extension not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"result": "failed to upload data",
			"error":  "Invalid file type. Only Excel sheets (xlsx, xlsm, xltm, xls) are allowed.",
		})
		return
	}

	tmpfile, err := os.CreateTemp("", "upload_*.xlsx")
	if err != nil {
		logger.Log.Debug("error while creating temp file ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"result": "failed to upload data",
			"error":  "failed to create temporary file",
		})
		return
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	logger.Log.Debug("uploading file ", tmpfile.Name())

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(tmpfile, file)
	if err != nil {
		logger.Log.Debug("error while copying content into temp file ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"result": "failed to upload data",
			"error":  "failed to copy contents into temp file",
		})
		return
	}

	fileName := strings.Split(tmpfile.Name(), `\`)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": "File uploaded successfully!  " + fileName[6],
		"error":  "",
	})
	go func() {
		var isUploaded bool = true
		if err = validateAndStoreExcel(tmpfile.Name()); err != nil {
			logger.Log.Debug("error while storing the file ", err.Error())
			isUploaded = false
		}
		uploadedFileStatusUpdate(fileName[6], isUploaded)
	}()
}

func isValidExcelFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".xlsx" || ext == ".xlsm" || ext == ".xltm" || ext == ".xls"
}

func validateAndStoreExcel(filename string) error {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		logger.Log.Debug("error while opening the file ", err.Error())
		return err
	}
	defer file.Close()

	// get the first sheetname by default
	sheetname := file.GetSheetName(0)
	rows, err := file.GetRows(sheetname)
	if err != nil || len(rows) < 1 {
		logger.Log.Debug("error while getting rows ", err)
		return errors.New("error while getting rows or no rows found in sheet 1")
	}

	//validate required headers of uploaded Excel sheet
	requiredHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "country", "postal", "phone", "email", "web"}
	for i := 0; i < len(requiredHeaders); i++ {
		if rows[0][i] != requiredHeaders[i] {
			logger.Log.Debug("missing required headers ", requiredHeaders[i])
			return errors.New("missing required headers in the excel sheet")
		}
	}

	// iterating over rows and creating employee model
	var employees []*models.Employee
	for i := 1; i < len(rows); i++ {
		employee := &models.Employee{
			FirstName: rows[i][0],
			LastName:  rows[i][1],
			Company:   rows[i][2],
			Address:   rows[i][3],
			City:      rows[i][4],
			Country:   rows[i][5],
			Postal:    rows[i][6],
			Phone:     rows[i][7],
			Email:     rows[i][8],
			Web:       rows[i][9],
		}
		employees = append(employees, employee)
	}

	// inserting into database
	err = insertIntoDataBase(employees)
	if err != nil {
		logger.Log.Debug("error while inserting employee ", err.Error())
		return errors.New("error while inserting employee in the database")
	}

	// inserting into redis with 5 min expiry
	insertIntoRedis(employees)
	return nil
}

func insertIntoDataBase(employees []*models.Employee) error {
	return dbs.InsertEmployee(employees)
}

func insertIntoRedis(employees []*models.Employee) {
	redisData, err := json.Marshal(employees)
	if err != nil {
		fmt.Println("error while marshalling data ", err.Error())
	}
	redisKey := Config.RedisConf.RedisUploadedDataKey
	redisKeyTtl := time.Duration(Config.RedisConf.RedisUploadedDataKeyTtl)
	if err := dbs.RClient.Set(redisKey, string(redisData), time.Minute*redisKeyTtl); err != nil {
		logger.Log.Debug("error while storing data into redis ", err.Error())
	}
}
