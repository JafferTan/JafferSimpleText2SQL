package main

import (
	"database/sql"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", "user=root password=Jaffer990206 host=localhost port=55003 dbname=pg sslmode=disable")
	//sqlStr := "select sum(s.amount) as 支出金额（元） from statement as s where s.month = '6月' and s.year = '2024'"
	sqlStr := "select * from statement"
	rows, _ := db.Query(sqlStr)
	defer rows.Close()
	//// 获取列名
	columns, _ := rows.Columns()
	//fmt.Println("表头列名:", columns)
	var records [][]string
	//// 根据列数创建值的切片
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	records = append(records, columns)
	for rows.Next() {
		// 读取一行数据到values切片中
		if err := rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
		}
		var row []string
		for _, v := range values {
			switch v := v.(type) {
			case string:
				row = append(row, v)
			case []byte:
				row = append(row, string(v))
			default:
				row = append(row, fmt.Sprintf("%v", v))
			}
		}
		records = append(records, row)
	}
	df := dataframe.LoadRecords(records)
	fmt.Println(df)
}
