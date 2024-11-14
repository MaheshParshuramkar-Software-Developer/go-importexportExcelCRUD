package models

import "time"

type Employee struct {
	Id        uint `gorm:"primary_key;auto_increment"`
	FirstName string
	LastName  string
	Company   string
	Address   string
	City      string
	Country   string
	Postal    string
	Phone     string
	Email     string
	Web       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ViewData struct {
	Id        uint
	FirstName string
	LastName  string
	Company   string
	Address   string
	City      string
	Country   string
	Postal    string
	Phone     string
	Email     string
	Web       string
}

type UpdateData struct {
	FirstName string `json:"first_name" binding:"omitempty"`
	LastName  string `json:"last_name" binding:"omitempty"`
	Company   string `json:"company" binding:"omitempty"`
	Address   string `json:"address" binding:"omitempty"`
	City      string `json:"city" binding:"omitempty"`
	Country   string `json:"country" binding:"omitempty"`
	Postal    string `json:"postal" binding:"omitempty"`
	Phone     string `json:"phone" binding:"omitempty"`
	Email     string `json:"email" binding:"omitempty"`
	Web       string `json:"web" binding:"omitempty"`
}

type UpdateDataUri struct {
	Id uint `uri:"id" binding:"required,number"`
}

type UploadedFileStatus struct {
	Filename  string    `json:"filename" binding:"required"`
	Status    string    `json:"status" binding:"omitempty"`
	CreatedAt time.Time `binding:"omitempty"`
	UpdatedAt time.Time `binding:"omitempty"`
}

//type FileUploadedStatus struct {
//	FileName string `json:"filename" binding:"required"`
//}
