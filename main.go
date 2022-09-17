package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type TodoList struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	List string             `json:"list,omitempty" bson:"list,omitempty"`
}

func CreateList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var todo TodoList
	_ = json.NewDecoder(request.Body).Decode(&todo)
	collection := client.Database("thepolyglotdeveloper").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, todo)
	json.NewEncoder(response).Encode(result)
}

func GetList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var todolist []TodoList
	collection := client.Database("thepolyglotdeveloper").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo TodoList
		cursor.Decode(&todo)
		todolist = append(todolist, todo)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(todolist)
}
func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateList).Methods("POST")
	router.HandleFunc("/get", GetList).Methods("GET")
	http.ListenAndServe(":3000", router)
}
