package stores

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	Summary string    `json:"summary"`
	Tags    []string  `json:"tags"`
	Done    bool      `json:"done"`
	Due     time.Time `json:"due"`
}

type TaskStore struct {
	Collection *mongo.Collection
}

func NewTaskStore(client *mongo.Client) *TaskStore {
	fmt.Println(DATABASE_NAME)
	return &TaskStore{
		Collection: client.Database(DATABASE_NAME).Collection("tasks"),
	}
}

func (ts *TaskStore) Create(task *Task) primitive.ObjectID {
	// add task to mongoDB
	result, err := ts.Collection.InsertOne(context.Background(), task)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	return result.InsertedID.(primitive.ObjectID)
}

func (ts *TaskStore) Get(id primitive.ObjectID) (*Task, error) {
	var task *Task

	filter := filterById(id)
	err := ts.Collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskStore) Delete(id primitive.ObjectID) error {
	filter := filterById(id)
	_, err := ts.Collection.DeleteOne(context.Background(), filter)
	return err
}

func (ts *TaskStore) Update(id primitive.ObjectID, task *Task) error {
	update := bson.D{{Key: "$set", Value: task}}
	_, err := ts.Collection.UpdateByID(context.Background(), id, update)
	fmt.Println(err)
	return err
}

func (ts *TaskStore) GetAll() ([]*Task, error) {
	var tasks []*Task
	cursor, err := ts.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func filterById(id primitive.ObjectID) bson.D {
	return bson.D{{Key: "_id", Value: id}}
}
