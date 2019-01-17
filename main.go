package main

import (
	"database/sql"
	"fmt"
	"github.com/flyleft/gprofile"
	_ "github.com/flyleft/gprofile"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strings"
)

type DataSource struct {
	Type     string `profile:"datasource.type"`
	Host     string `profile:"datasource.host"`
	Port     string `profile:"datasource.port"`
	Username string `profile:"datasource.username"`
	Password string `profile:"datasource.password"`
	Database string `profile:"datasource.database"`
}

type DataSourcePath struct {
	Path string `profile:"datasource.path"`
	Type string `profile:"datasource.type"`
}

var Path *DataSourcePath

func main() {
	//fmt.Println("hello")
	//var a = 5
	//var b = 2
	//fmt.Println(add(a, b))
	//useDb()
	//getEnv()
	//conf := getConfig()
	//fmt.Println(conf.Type)
	//fmt.Println(conf)
	//path := getDBPath()
	//fmt.Println(path)
	//http.HandleFunc("/", reqHandlerTest)	//设置访问的路由
	Path = getDBPath()
	http.HandleFunc("/usedb", testRequestDB)
	err := http.ListenAndServe(":8989", nil) //设置监听的端口
	checkErr(err)
}

func testRequestDB(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("path", req.URL.Path)
	useDb()
}

func reqHandlerTest(res http.ResponseWriter, req *http.Request) {
	req.ParseForm() //解析参数，默认是不会解析的
	fmt.Println(req.Form)
	fmt.Println("path", req.URL.Path)
	fmt.Println("scheme", req.URL.Scheme)
	fmt.Println(req.Form["url_long"])
	for k, v := range req.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(res, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func getDBPath() *DataSourcePath {
	env, err := gprofile.Profile(&DataSourcePath{}, "config.yml", true)
	checkErr(err)
	return env.(*DataSourcePath)
}

func getConfig() *DataSource {
	env, err := gprofile.Profile(&DataSource{}, "config.yml", true)
	checkErr(err)
	return env.(*DataSource)
}

func add(a, b int) int {
	return a + b
}

func getEnv() {
	a := os.Getenv("a")
	fmt.Println(a)
}

func useDb() {
	/*DSN数据源名称
	[username[:password]@][protocol[(address)]]/dbname[?param1=value1¶mN=valueN]
	user@unix(/path/to/socket)/dbname
	user:password@tcp(localhost:5555)/dbname?charset=utf8&autocommit=true
	user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname?charset=utf8mb4,utf8
	user:password@/dbname
	无数据库: user:password@/
	*/
	db, err := sql.Open(Path.Type, Path.Path)
	checkErr(err)
	//db.Query("drop database if exists test-mysql")
	//db.Query("create database test-mysql")
	//db.Query("use test-mysql")
	//db.Query("create table test_db(c1 int, c2 varchar(20), c3 varchar(20))")
	db.Query("truncate test_db")
	db.Query("insert into test_db values (101, '姓名1', 'address1'), (102, '姓名2', 'address2'), (103, '姓名3', 'address3'), (104, '姓名4', 'address4')")
	query, err := db.Query("select * from test_db")
	checkErr(err)
	//v := reflect.ValueOf(query)
	//fmt.Println(v)
	printResult(query)
	db.Close()
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}

func printResult(query *sql.Rows) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	for k, v := range results { //查询出来的数组
		fmt.Println(k, v)
	}
}
