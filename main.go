package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

//发布言论、评论、点赞、回复

type Liuyan struct {
	fid         int    `DB:"fid"` //被评论id，fid=0则该留言为父留言，其不是别的留言的评论
	id          int    `DB:"id"`
	userName    string `DB:"userName"`
	content     string `DB:"content"`
	time        string `DB:"time"`
	num_zan     int    `DB:"num_zan"`
	num_conment int    `DB:"num_conment"`
}

func main() {

	DB := initDB()

	//zanDB(DB,1)

	r := gin.Default()

	//获取留言
	//getliuyan/xxxx留言父id
	r.GET("/getLiuyan/:id", func(c *gin.Context) {
		s := c.Param("id")
		id, _ := strconv.Atoi(s)
		//selectDB(DB,id)
		rows, _ := DB.Query("SELECT * FROM liuyan WHERE fid =?;", id)

		for rows.Next() {
			var s Liuyan
			_ = rows.Scan(&s.fid, &s.id, &s.userName, &s.content, &s.time, &s.num_zan, &s.num_conment)

			c.JSON(200, gin.H{
				"fid":         s.fid,
				"id":          s.id,
				"userName":    s.userName,
				"content":     s.content,
				"time":        s.time,
				"num_zan":     s.num_zan,
				"num_conment": s.num_conment,
			})
		}
		rows.Close()
	})

	//点赞
	r.POST("/zan", func(c *gin.Context) {
		s := c.PostForm("id")
		id, _ := strconv.Atoi(s)

		zanDB(DB, id) //点赞数加一

		c.JSON(http.StatusOK, gin.H{
			"status": "SUCCESS",
			"id":     id,
		})
	})

	//删除评论或留言
	r.POST("/delete", func(c *gin.Context) {
		s := c.PostForm("id")
		id, _ := strconv.Atoi(s)

		deleteDB(DB, id) //点赞数加一

		c.JSON(http.StatusOK, gin.H{
			"status": "DELETE SUCCESS",
			"id":     id,
		})
	})


	//发布留言
	r.POST("/postLiuyan", func(c *gin.Context) {
		pUserName := c.PostForm("username")
		pContent := c.PostForm("content")

		var s Liuyan
		s.fid = 0
		s.num_zan = 0
		s.num_conment = 0
		s.userName = pUserName
		s.content = pContent
		s.time = getTime()
		row := DB.QueryRow("SELECT max(id) from liuyan")
		var n int
		row.Scan(&n)
		n++
		s.id = n

		insertDB(DB, s)

		c.JSON(200, gin.H{
			"fid":         s.fid,
			"id":          s.id,
			"userName":    s.userName,
			"content":     s.content,
			"time":        s.time,
			"num_zan":     s.num_zan,
			"num_conment": s.num_conment,
		})
	})

	//发布评论/回复
	r.POST("/postComment", func(c *gin.Context) {
		sid := c.PostForm("fid")
		pFid, _ := strconv.Atoi(sid)
		pUserName := c.PostForm("username")
		pContent := c.PostForm("content")

		var s Liuyan
		s.fid = pFid
		s.num_zan = 0
		s.num_conment = 0
		s.userName = pUserName
		s.content = pContent
		s.time = getTime()
		row := DB.QueryRow("SELECT max(id) from liuyan")
		var n int
		row.Scan(&n)
		n++
		s.id = n

		insertDB(DB, s)
		pinglunDB(DB,pFid)

		c.JSON(200, gin.H{
			"fid":         s.fid,
			"id":          s.id,
			"userName":    s.userName,
			"content":     s.content,
			"time":        s.time,
			"num_zan":     s.num_zan,
			"num_conment": s.num_conment,
		})
	})



	r.Run()
}

func initDB() *sql.DB {
	DB, _ := sql.Open("mysql", "root:Lmqzuishuai01@tcp(localhost:3306)/test?charset=utf8")

	//DB.SetConnMaxLifetime(100)

	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return &sql.DB{}
	}
	fmt.Println("DB connect success")
	return DB
}

func insertDB(db *sql.DB, m Liuyan) {//增

	_, err := db.Exec("INSERT INTO liuyan VALUES (?,?,?,?,?,?,?)", m.fid, m.id, m.userName, m.content, m.time, 0, 0)

	if err != nil {
		panic(err)
	}
}

func zanDB(db *sql.DB, m int) {
	rows, _ := db.Query("SELECT * FROM liuyan WHERE id =?;", m)
	var a int
	for rows.Next() {
		var s Liuyan
		_ = rows.Scan(&s.fid, &s.id, &s.userName, &s.content, &s.time, &a, &s.num_conment)
	}
	a++
	stmt, err := db.Exec("UPDATE liuyan SET num_zan=? WHERE id=?", a, m)
	if err != nil {
		log.Fatal(err)
	}
	stmt.RowsAffected()
}

func deleteDB(db *sql.DB,m int)  {
	stmt,_:=db.Exec("DELETE from liuyan WHERE id =?",m)
	stmt,_=db.Exec("DELETE from liuyan WHERE fid =?",m)
	stmt.RowsAffected()
}

func pinglunDB(db *sql.DB, m int) {
	rows, _ := db.Query("SELECT num_zan FROM liuyan WHERE id =?;", m)
	var a int
	for rows.Next() {
		var s Liuyan
		_ = rows.Scan(&s.fid, &s.id, &s.userName, &s.content, &s.time, &s.num_zan, &a)
	}
	a++
	stmt, err := db.Exec("UPDATE liuyan SET num_conment=? WHERE id=?", a, m)
	if err != nil {
		log.Fatal(err)
	}
	stmt.RowsAffected()
}

func selectALL(db *sql.DB) {//展示所有数据

	rows, _ := db.Query("SELECT * FROM liuyan;")

	for rows.Next() {
		var s Liuyan
		_ = rows.Scan(&s.fid, &s.id, &s.userName, &s.content, &s.time, &s.num_zan, &s.num_conment)
		fmt.Println(s)
	}
	rows.Close()
}

func selectDB(db *sql.DB, fid int) { //根据被评论id查找评论

	rows, _ := db.Query("SELECT * FROM liuyan WHERE fid =?;", fid)

	for rows.Next() {
		var s Liuyan
		_ = rows.Scan(&s.fid, &s.id, &s.userName, &s.content, &s.time, &s.num_zan, &s.num_conment)
		fmt.Println(s)
	}
	rows.Close()
}

func getTime() string {
	now := time.Now()
	//fmt.Println(now)
	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	//fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}
