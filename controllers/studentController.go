package controllers

import (
	"assignment-project/database"
	"assignment-project/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatedStudent(ctx *gin.Context) {
	db := database.GetDB()

	var newStudent models.Student

	if err := ctx.BindJSON(&newStudent); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	student := models.Student{
		Name : newStudent.Name,
		Age : newStudent.Age,
		Scores : newStudent.Scores,
	}

	err := db.Create(&student).Error

	if err != nil {
		fmt.Println("Error while creating new student : ", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success" : true,
		"data" : student,
	})

}

func GetAllStudent(ctx *gin.Context) {
	db := database.GetDB()
	var allStudent = []models.Student{}

	res := db.Preload("Scores").Find(&allStudent)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    allStudent,
	})
}

func UpdateOrder(ctx *gin.Context){
	id:= ctx.Param("ID")
	db:= database.GetDB()
	var updateStudent models.Student 

	if err := ctx.BindJSON(&updateStudent); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	studentID, err := strconv.Atoi(id)
	if err !=nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}

	student := models.Student{
		Name : updateStudent.Name,
		Age : updateStudent.Age,
		Scores : updateStudent.Scores,
	}

	previousStudent := models.Student{}
	result := db.Preload("Scores").First(&previousStudent, uint(studentID))
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to retrieve student"})
		return
	}

	trx :=db.Begin()

	previousStudent.Name = student.Name
	previousStudent.Age = student.Age

	for i, newStudent := range student.Scores{
		if i < len(previousStudent.Scores){
			previousStudent.Scores[i].AssignmentTitle= newStudent.AssignmentTitle
			previousStudent.Scores[i].Description= newStudent.Description
			previousStudent.Scores[i].Score= newStudent.Score
			if err := trx.Save(&previousStudent.Scores[i]).Error; err != nil {
				trx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error":"failed to updated scores"})
				return
			}
		}else{
			previousStudent.Scores = append(previousStudent.Scores, newStudent)
			if err := trx.Create(&previousStudent.Scores[i]).Error; err !=nil{
				trx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create new scores"})
				return
			}
		}
	}

	if err:= trx.Save(&previousStudent).Error; err !=nil{
		trx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to update student"})
		return
	}

	trx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"error": true,
		"data" : previousStudent,
	})

}

func DeleteStudent(ctx *gin.Context){
	id := ctx.Param("ID")
	db := database.GetDB()

	var deleteStudent models.Student 

	studentID, err := strconv.Atoi(id)
	if err !=nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	result := db.Preload("Scores").First(&deleteStudent, studentID)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to find student"})
		return
	}

	trx := db.Begin()

	if err := trx.Delete(&deleteStudent.Scores).Error; err != nil{
		trx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to remove score student"})
		return
	}

	if err := trx.Delete(&deleteStudent).Error; err != nil{
		trx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to remove student"})
		return
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"success":true,
		"data": nil,
	})

}