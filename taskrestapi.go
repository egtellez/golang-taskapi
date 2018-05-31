package main

import (
    "github.com/gin-gonic/gin"
	"io"
    "io/ioutil"
	"os"
	"log"
	"encoding/json"
	"strconv"
)

type Task struct {
		Id			  int "json:'Id'"
		Name          string "json:'Name'"
		Description   string "json:'Description'"
		Owner		  string "json:'Owner'"
		State	  	  string "json:'DeadLine'"
		Priority	  int "json:Priority"	
}

func main() {

	f, _ := os.Create("gin.log")
    //gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()
	
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(200, getTasks())
	})
	
	r.POST("/tasks", func(c *gin.Context) {
	    
		var err error
		c.Header("Content-Type", "application/json; charset=utf-8")
		t:= Task{}
		if err = c.BindJSON(&t); err != nil {
			c.JSON(400, gin.H{
				"error":  "json decoding : " + err.Error(),
				"status": 400,
		})
		return
		}else{
			result := saveTask(t);
			t.State = "New"
			c.JSON(200, gin.H{
				"success" : result,
			})
		}
		return
	})
	
	r.DELETE("tasks/:taskId", func(c *gin.Context){
		var taskid int64
		taskid, err := strconv.ParseInt(c.Param("taskId"), 10, 32)
		if(err!=nil){
			c.JSON(400, gin.H{
				"message" : "taskId not an int",
			})
		} else{
		result := deleteTask(taskid)
		c.JSON(200, gin.H{
				"success" : result,
		})
		}
	})
	
	r.Run() // listen and serve on 0.0.0.0:8080
}

func getTasks() (tasks []Task) {
   tasks = readTasksFile()
   return tasks
}

func readTasksFile() (result []Task){
	fileName := "tasks.json"
	file, err := ioutil.ReadFile(fileName);
	taskDB := []Task{}
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		if err = json.Unmarshal([]byte(file), &taskDB); err != nil {
			log.Fatal(err)
			return nil
		}
	}
	return taskDB
}

func saveTask(task Task) (success bool){
	
	taskDB := readTasksFile()
	taskDB = append(taskDB, task)
	
	success = writeTasksToFile(taskDB)
	return success
}

func deleteTask(taskId int64) (success bool){
	task:= Task{}
	task.Id = int(taskId)
	taskDB := readTasksFile()
	taskDB2 := []Task{}
	
	for i := range taskDB {
		if taskDB[i].Id != task.Id {
			taskDB2 = append(taskDB2, taskDB[i])
		}
	}
	success = writeTasksToFile(taskDB2)
	return success
}

func writeTasksToFile(taskDB []Task) (success bool){
	fileName := "tasks.json"
	d1, err := json.Marshal(taskDB)
	if err = ioutil.WriteFile(fileName, d1, 0644); err!=nil {
		log.Fatal(err)
		return false
	} else {
		log.Println(taskDB)
		return true
	}
}