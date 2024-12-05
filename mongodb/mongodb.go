package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Scene struct {
	SceneDesc string `bson:"scene_desc"`
	//Prompt         string `bson:"prompt"`
	SuperSceneCode string `bson:"super_scene_code"`
	SceneCode      string `bson:"scene_code"`
}

func FindLevel1SceneCode(question string) string {
	//这里要开始找一层场景
	level1Scenes := findLevel1ScenesFromMongoDB()
	jsonData, _ := json.Marshal(level1Scenes)
	level1ScenesString := string(jsonData)
	content, _ := ioutil.ReadFile("prompt_file/prompt.txt")
	contentStr := string(content)
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	contentStr = strings.Replace(contentStr, "{scenarios}", level1ScenesString, -1)
	contentStr = strings.Replace(contentStr, "{user_input}", question, -1)
	contentStr = strings.Replace(contentStr, "{now}", formattedTime, -1)
	return contentStr
}

func FindPromptBySceneCode(sceneCode string) string {
	uri := "mongodb://localhost:55002"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	DB := client.Database("local")
	table := DB.Collection("prompt_management")
	filter := bson.D{{"scene_code", bson.D{{"$eq", sceneCode}}}} // 查询年龄小于等于3的，这里特别有意思，能够使用$lte这种方法，类似这样的，MongoDB还提供了很多其他的查询方法，比如$gt等等
	ctx := context.Background()
	cursor, _ := table.Find(ctx, filter)
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		byteData, _ := bson.Marshal(result)
		return string(byteData)
	}
	return ""
}

func findLevel1ScenesFromMongoDB() []*Scene {
	uri := "mongodb://localhost:55002"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	DB := client.Database("local")
	table := DB.Collection("prompt_management")
	filter := bson.D{{"super_scene_code", bson.D{{"$eq", nil}}}} // 查询年龄小于等于3的，这里特别有意思，能够使用$lte这种方法，类似这样的，MongoDB还提供了很多其他的查询方法，比如$gt等等
	ctx := context.Background()
	cursor, err := table.Find(ctx, filter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	var res []*Scene
	for _, result := range results {
		var scene *Scene
		bytedata, _ := bson.Marshal(result)
		err := bson.Unmarshal(bytedata, &scene)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		res = append(res, scene)
	}
	return res
}
