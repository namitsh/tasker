package tasker

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// mongoDB Configuration

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("tasker").Collection("tasks")
}

// model

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

// methods

func CreateTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}

func GetAllTasks() ([]*Task, error) {
	filter := bson.D{{}}
	return filterTasks(filter)
}

func PendingTasks() ([]*Task, error) {
	filter := bson.M{"completed": false}
	return filterTasks(filter)
}

func FinishedTasks() ([]*Task, error) {
	filter := bson.M{"completed": true}
	return filterTasks(filter)
}

func DeleteTask(text string) error {
	filter := bson.M{"text": text}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no tasks were deleted")
	}
	return nil
}

func filterTasks(filter interface{}) ([]*Task, error) {
	var tasks []*Task
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	for cur.Next(ctx) {
		var task Task
		if err := cur.Decode(&task); err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	if err := cur.Err(); err != nil {
		return tasks, err
	}
	defer cur.Close(ctx)
	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}
	return tasks, nil
}

func MarkCompleted(text string) error {
	var result Task
	filter := bson.M{"text": text}
	update := bson.M{"$set": bson.M{"completed": true}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}
