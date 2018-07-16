package main

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()
var projectID = "secret-primacy-210308"
var client, clientErr = datastore.NewClient(ctx, projectID)

type Task struct {
	Value   string
	Created time.Time
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if clientErr != nil {
		log.Fatalf("Failed to create client: %v", clientErr)
	}
	inputValue := r.FormValue("input")
	if len(inputValue) != 0 {
		kind := "Task"
		key := datastore.NameKey(kind, "", nil)
		task := Task{
			Value:   inputValue,
			Created: time.Now(),
		}
		if _, err := client.Put(ctx, key, &task); err != nil {
			log.Fatalf("Failed to save task: %v", err)
		}
		fmt.Fprintf(w, "Saved %v to kind %s", task.Value, kind)
	}
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	if clientErr != nil {
		log.Fatalf("Failed to create client: %v", clientErr)
	}
	var tasks []*Task
	query := datastore.NewQuery("Task").Order("Created")
	_, err := client.GetAll(ctx, query, &tasks)
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	for _, v := range tasks {
		fmt.Fprintf(w, "%v\n", (*v).Value)
	}
}

func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/retrieve", retrieveHandler)
	http.ListenAndServe(":8080", nil)
}
