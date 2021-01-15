package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	r := gin.Default()

	r.StaticFS("/static", http.Dir("./static"))

	r.GET("/test", func(c *gin.Context) {
		u := c.DefaultQuery("url", "")
		sql := c.DefaultQuery("sql", "")
		fmt.Println(u, sql)
		resp, err := http.Get(u + "/sql?sql=" + url.QueryEscape(sql))
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body),err)
		c.JSON(200, gin.H{
			"data":  string(body),
			 "code":200,
		})
	})
	r.Run()
	//resp, err := http.Get("http://localhost:6001/sql?sql=select%20*%20from%20%20userinfo")
	////resp, err := http.Get("http://localhost:6000/sql?sql=insert%20into%20userinfo(name)%20values%20('666')")
	////resp, err := http.Get("http://localhost:6000/sql?sql=CREATE%20TABLE%20userinfo(name%20TEXT)")
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body),err)
}
