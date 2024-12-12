package main

import (
	"JafferSimpleText2SQL/db"
	docs "JafferSimpleText2SQL/docs"
	"JafferSimpleText2SQL/elasticsearch"
	llm "JafferSimpleText2SQL/model"
	"JafferSimpleText2SQL/mongodb"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger scenarios
	"net/http"
	"strings"
	"time"
)

type UserInput struct {
	Question string `json:"question" form:"question" binding:"required"`
}

type SceneRecognition struct {
	Scenarios string `json:"scenarios"`
	Thinking  string `json:"thinking"`
}

type Text2SQLResponse struct {
	SQL       string     `json:"sql"`
	DataFrame [][]string `json:"dataframe"`
	DimInfo   DimInfo    `json:"diminfo"`
}

type DimInfo struct {
	From    string              `json:"from"`
	Filters []map[string]string `json:"filters"`
}

//	@BasePath	/api/v1

// Text2SQL godoc
//
// @Summary Convert text to SQL query
// @Schemes http
// @Description Accepts user input in JSON format, binds it to the UserInput struct, and attempts to convert it to a relevant SQL query. If there is an error in parsing the JSON request body, it will be logged.
// @Tags Text2SQL
// @Param userInput body UserInput true "The user input containing details for generating the SQL query. See the UserInput struct for details."
// @Accept json
// @Produce json
// @Success		200	{string}	Helloworld
// @Router /example/Text2SQL [post]
func Text2SQL(g *gin.Context) {
	var userInput UserInput
	g.ShouldBindJSON(&userInput)
	res := ChatWithGPT(userInput.Question)
	//1.先进行场景识别，将当前数据的场景获取得到，然后再进行下一步处理
	var scencesRecognition SceneRecognition
	json.Unmarshal([]byte(res), &scencesRecognition)
	prompt := mongodb.FindPromptBySceneCode(scencesRecognition.Scenarios)
	//2.将prompt进行注入,rag信息
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	prompt = strings.ReplaceAll(prompt, "{user_input}", userInput.Question)
	ragInfoList := elasticsearch.FindRagInfoBySceneCode(scencesRecognition.Scenarios, userInput.Question)
	jsonStr, _ := json.Marshal(ragInfoList)
	prompt = strings.ReplaceAll(prompt, "{esDims}", string(jsonStr))
	prompt = strings.ReplaceAll(prompt, "{now}", formattedTime)
	//3.生成SQL
	text2SQL := ChatWithGPT(prompt)
	//fmt.Println(text2SQL)
	var response Text2SQLResponse
	json.Unmarshal([]byte(text2SQL), &response)
	df := db.ReadDataFrameFromPG(response.SQL)
	response.DataFrame = df.Records()
	//fmt.Println(df)
	g.JSON(http.StatusOK, response)
}

func ChatWithGPT(question string) string {
	//os.Setenv("OPENAI_API_KEY", "sk-svcacct-z1y21xwzywlV-3xjRvI9QOd3C1vYVDJMjZRHj5EgME7z4Hp1cn3YqeGEI46CBT3BlbkFJGbu_tg5uKRFEQ7v1dNWB06o-wzZSHOxUMbOanuc5D31tK5e4JQ_2xStxq_OjwA")
	//llm, err := openai.New(openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"))
	prompt := mongodb.FindLevel1SceneCode(question)
	completion, _ := llm.CallLLM(prompt)
	return completion.Choices[0].Message.Content
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.POST("/Text2SQL", Text2SQL)
		}

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run("0.0.0.0:8080")

}
