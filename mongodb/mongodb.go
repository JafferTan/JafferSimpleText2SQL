package mongodb

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic/encoder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

func FindLevel1SceneCode(question string) string {
	//这里要开始找一层场景
	findLevel1ScenesFromMongoDB()
	content, _ := ioutil.ReadFile("prompt_file/prompt.txt")
	//fmt.Println(string(content))
	return string(content)
}

func findLevel1ScenesFromMongoDB() {
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
	for _, result := range results {
		v, err := encoder.Encode(result, encoder.SortMapKeys)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println(string(v))
	}
}

func main() {
	uri := "mongodb://localhost:55002"
	clientOptions := options.Client().ApplyURI(uri)
	_, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)

	}
	content, _ := ioutil.ReadFile("prompt_file/prompt.txt")
	fmt.Println(string(content))

}
