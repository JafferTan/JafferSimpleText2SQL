package main

import (
	docs "JafferSimpleText2SQL/docs"
	llm "JafferSimpleText2SQL/model"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/tmc/langchaingo/llms"
	"net/http"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	res := ChatWithGPT()
	g.JSON(http.StatusOK, res)
}

func ChatWithGPT() string {
	//os.Setenv("OPENAI_API_KEY", "sk-svcacct-z1y21xwzywlV-3xjRvI9QOd3C1vYVDJMjZRHj5EgME7z4Hp1cn3YqeGEI46CBT3BlbkFJGbu_tg5uKRFEQ7v1dNWB06o-wzZSHOxUMbOanuc5D31tK5e4JQ_2xStxq_OjwA")
	//llm, err := openai.New(openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"))
	prompt := "What would be a good company name for a company that makes colorful socks?"
	completion, _ := llm.CallLLM(prompt)
	return completion.Choices[0].Message.Content
}

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")

	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")

}
