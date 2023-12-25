package main

import (
	"database/sql"
	"flag"
	"fmt"
	mpset "github.com/deckarep/golang-set"
	_ "github.com/godror/godror"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

type Config struct {
	Login struct {
		Address  string `yml:"address"`
		Username string `yml:"username"`
		Password string `yml:"password"`
	} `yml:"login"`
	Check []map[string][]string `yml:"check"`
}

var (
	config   Config
	querySql string = `select index_name,column_name,column_position from all_ind_columns where table_name = upper('%s') and column_position = %d and column_name = upper('%s')`
)

// 登陆数据库
func loginDb() *sql.DB {
	dsn := "user=" + config.Login.Username + " password=" +
		config.Login.Password + " connectString=" +
		//config.Login.Address + " timezone=" + "Asia/Shanghai"
		config.Login.Address + " timezone=" + "local"
	//fmt.Println(dsn)
	client, err := sql.Open("godror", dsn)
	handleErr(err)

	return client
}

// 索引字段比较
func compareIdxField(query string, client *sql.DB) bool {
	//query
	rows, err := client.Query(query)
	handleErr(err)
	// 查询结果判断
	if rows.Next() {
		// 打印查询结果
		var indexName, columnName, columnPosition string
		err = rows.Scan(&indexName, &columnName, &columnPosition)
		handleErr(err)
		log.Debugf("+++ indexName -> %s ;columnName -> %s ;columnPosition -> %s", indexName, columnName, columnPosition)
		//fmt.Println("pass")
		for rows.Next() {
			err = rows.Scan(&indexName, &columnName, &columnPosition)
			handleErr(err)
			log.Debugf("+++ indexName -> %s ;columnName -> %s ;columnPosition -> %s", indexName, columnName, columnPosition)
		}
		return true
	} else {
		//fmt.Println("no pass")
		return false
	}
}

// 执行检查
func doCheck(client *sql.DB) {
	for _, check := range config.Check {
		for tab, v := range check {
			for _, vv := range v {
				arrField := strings.Split(vv, ",")
				// single field 单个字段索引判断
				if len(arrField) == 1 {
					log.Debug("单索引检测")
					//fmt.Println("single field: ", tab, vv)
					query := fmt.Sprintf(querySql, tab, 1, arrField[0])
					//fmt.Println(query)
					log.Debug(query)
					flag := compareIdxField(query, client)

					// 执行sql结果判断
					if flag {
						log.Infof("=== pass single field: tab -> %s ;idx -> %s ", tab, vv)
					} else {
						log.Infof("--- false single field: tab -> %s ;idx -> %s ", tab, vv)
					}
					// mulite field
				} else {
					log.Debug("组合索引检测：逐个字段比对")
					//fmt.Println("mulite field: ", tab, vv)

					// 组合索引,逐个判断是否独立索引
					flags := mpset.NewSet()

					for _, vField := range arrField {
						//fmt.Println("+++", tab, vField)
						// 单独索引
						query := fmt.Sprintf(querySql, tab, 1, vField)
						log.Debug(query)
						flag := compareIdxField(query, client)
						flags.Add(flag)
					}
					// 判断是否继续
					if flags.Cardinality() == 1 {
						log.Infof("=== pass multi field by all single field: tab -> %s ;idx -> %s", tab, vv)
						// 打印空行
						log.Info(strings.Repeat("-", 20))
						continue
					}
					flags.Clear()

					// 组合索引
					log.Debug("组合索引检测：组合顺序比对")
					for i, vField := range arrField {

						query := fmt.Sprintf(querySql, tab, i+1, vField)
						log.Debug(query)
						flag := compareIdxField(query, client)
						flags.Add(flag)
					}
					// 对组合索引结果进行判断
					if flags.Contains(false) {
						log.Infof("--- false multi field by multi field: tab -> %s ;idx -> %s", tab, vv)
					} else {
						log.Infof("=== pass multi field by multi field: tab -> %s ;idx ->%s", tab, vv)
					}
				}
				// 打印空行
				log.Info(strings.Repeat("-", 20))
			}
		}
	}
}

// 读取yml
func readConfig() {
	bytes, err := os.ReadFile("config.yml")
	handleErr(err)
	// 解码
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
}

// 设置默认日志级别
func init() {
	level := flag.String("level", "info", "log level, value: debug|info")
	flag.Parse()
	// 初始化日志级别
	switch *level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	f, err := os.OpenFile("check.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	handleErr(err)
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func main() {
	readConfig()
	client := loginDb()
	defer client.Close()
	doCheck(client)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
