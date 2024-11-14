package dbs

import (
	"errors"

	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb(mySqlConf *models.MySqlConfig) {
	var err error
	dsn := mySqlConf.UserName + ":" + mySqlConf.Password + "@tcp(" + mySqlConf.Host + ":" + mySqlConf.Port + ")/" + mySqlConf.Database + "?parseTime=true&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database ", err.Error())
	}
}

func SyncDb() {
	if err := DB.AutoMigrate(&models.Employee{}, &models.UploadedFileStatus{}); err != nil {
		logger.Log.Fatal("error while migrating database ", err.Error())
	}
}

func InsertEmployee(employee []*models.Employee) error {
	response := DB.Create(employee)
	if response.Error != nil {
		logger.Log.Debug("error while inserting into database ", response.Error.Error())
		return errors.New("error while uploading records")
	}
	return nil
}

func GetAllEmployee() ([]models.ViewData, error) {
	var employees []models.ViewData
	response := DB.Model(&models.Employee{}).Find(&employees)
	return employees, response.Error
}

func GetEmployee(id uint) (models.Employee, error) {
	var employee models.Employee
	employee.Id = id
	response := DB.First(&employee)
	return employee, response.Error
}

func UpdateEmployee(updateData *models.Employee) error {
	result := DB.Model(&updateData).Updates(updateData)
	if result.Error != nil {
		logger.Log.Debug("error while updating records", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records to update")
	}
	return nil
}

func DeleteEmployee(id uint) error {
	result := DB.Delete(&models.Employee{}, id)
	if result.Error != nil {
		logger.Log.Debug("error while deleting records", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records to delete")
	}
	return nil
}

func InsertFileStatus(uploadedFile *models.UploadedFileStatus) error {
	result := DB.Create(&uploadedFile)
	if result.Error != nil {
		logger.Log.Debug("error while inserting records", result.Error)
		return result.Error
	}
	return nil
}

func GetUploadedFileStatus(filename string) (models.UploadedFileStatus, error) {
	var uploadedFile models.UploadedFileStatus
	uploadedFile.Filename = filename
	result := DB.Where(&uploadedFile).First(&uploadedFile)
	return uploadedFile, result.Error
}
