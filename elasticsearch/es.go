package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"strings"
)

type SceneRagInfo struct {
	SceneCode string `json:"scene_code"`
	Value     string `json:"value"`
}

func convertResultToSceneRagInfo(result map[string]interface{}) ([]SceneRagInfo, error) {
	var sceneRagInfos []SceneRagInfo

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]

		// 将_source中的数据转换为字节数组，以便后续用json.Unmarshal解析
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var sceneRagInfo SceneRagInfo
		err = json.Unmarshal(sourceBytes, &sceneRagInfo)
		if err != nil {
			return nil, err
		}

		sceneRagInfos = append(sceneRagInfos, sceneRagInfo)
	}

	return sceneRagInfos, nil
}

func FindRagInfoBySceneCode(sceneCode, question string) []SceneRagInfo {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:55000"},
		Username:  "elastic",
		Password:  "5fVNSR06Uw9nQMOBeGg6",
	}
	es, _ := elasticsearch.NewClient(cfg)
	filter := `{"query":{"bool":{"must":[{"term":{"scene_code":"{scene_code}"}},{"match":{"value":"{question}"}}]}}}`
	filter = strings.ReplaceAll(filter, "{scene_code}", sceneCode)
	filter = strings.ReplaceAll(filter, "{question}", question)
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("elasticsearch"),
		es.Search.WithBody(strings.NewReader(filter)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding search results: %v", err)
	}

	// 将查询结果转换为SceneRagInfo结构体数组
	sceneRagInfos, err := convertResultToSceneRagInfo(result)
	if err != nil {
		log.Fatalf("Error converting search results to SceneRagInfo: %v", err)
	}
	return sceneRagInfos
}

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:55000"},
		Username:  "elastic",
		Password:  "5fVNSR06Uw9nQMOBeGg6",
	}
	es, _ := elasticsearch.NewClient(cfg)
	filter := `{"query":{"bool":{"must":[{"term":{"scene_code":"{scene_code}"}},{"match":{"value":"{question}"}}]}}}`
	sceneCode := "TILE_DEFECT"
	question := "白色"
	filter = strings.ReplaceAll(filter, "{scene_code}", sceneCode)
	filter = strings.ReplaceAll(filter, "{question}", question)
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("elasticsearch"),
		es.Search.WithBody(strings.NewReader(filter)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding search results: %v", err)
	}

	// 将查询结果转换为SceneRagInfo结构体数组
	sceneRagInfos, err := convertResultToSceneRagInfo(result)
	if err != nil {
		log.Fatalf("Error converting search results to SceneRagInfo: %v", err)
	}

	// 输出转换后的结构体数组内容
	for _, sceneRagInfo := range sceneRagInfos {
		fmt.Printf("SceneCode: %s, Value: %s\n", sceneRagInfo.SceneCode, sceneRagInfo.Value)
	}

}
