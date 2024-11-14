// README.md
# Project Choice Golang Assignment

### Introduction
Project Choice Golang Assignment is a Golang program to import data from an Excel file, store it into MySQL, and cache the data in Redis. Create a simple CRUD (Create, Read, Update, Delete) system to view imported data, edit records, and update the changes to
both the database and cache using a Gin Framework .

### Project Choice Golang Assignment Features
* Users can Upload Excel file which contains data employees data.
* Users can View the data which was uploaded through Excel file.
* Users can Update the data.
* Users can Delete the data.
* Users can check the status of file uploaded.

### Prerequisites
Before running the application, ensure you have the following installed
* Go: Version 1.20 or later
* MySQL: Version 8.0 or later
* Redis: Version 7.4 or later

### Installation Guide
* Clone the repository.
* Run go mod tidy to install all dependencies.
* Use your locally installed mysql, redis, and Go.

### Usage
* Run go run . or go run main.go to start the application.
* Connect to the API using Postman on port 8080.

### API Endpoints
| HTTP Verbs | Endpoints               | Action                                      |
|------------|-------------------------|---------------------------------------------|
| GET        | /api/ping               | To check the status of Server               |
| POST       | /api/uploadData         | To upload the Excel File                    |
| GET        | /api/viewData           | To View the Uploaded Data                   |
| PUT        | /api/updateData/:id     | To Update the data by Id                    |
| DELETE     | /api/deleteData/:id     | To Delete the data by Id                    |
| POST       | /api/fileUploadedStatus | To Get the uploaded file status by filename |


### Technologies Used
* [Golang](https://go.dev/) This is an open-source programming language supported by Google, Easy to learn and great for teams ,Built-in concurrency and a robust standard library ,Large ecosystem of partners, communities, and tools.
* [Gin](https://gin-gonic.com/) A high-performance HTTP web framework for Go, used for creating the API.
* [MySql](https://www.mysql.com/) A relational database management system for persistent data storage.
* [Redis](https://www.googleadservices.com/pagead/aclk?sa=L&ai=DChcSEwjivv2bv9OJAxXMEXsHHSq3K0gYABABGgJ0bQ&co=1&ase=2&gclid=Cj0KCQiA0MG5BhD1ARIsAEcZtwTpS7tYG7cuiwbm_Leu9svCUflAuXzGsX0sAJ7w3233ww5ivImCkPwaAgAoEALw_wcB&ohost=www.google.com&cid=CAESVeD24vXR3uzuMD1IOwDmLxVk7K8I5V80_DhGvo8y6s6sXYKDjDqeqmwPWh-mXXE5C9W49Q-Mh8WEKH8hX-LkQ4oy7RFjb41kzr-oM0HujDWqdM68W5k&sig=AOD64_3CLaszBKyGGgH4yL_tFYm0gUPSGw&q&nis=4&adurl&ved=2ahUKEwi75vabv9OJAxXFdfUHHd3dAkMQ0Qx6BAgLEAE) An in-memory data structure store used for caching.

### Authors
* Mahesh Parshuramkar