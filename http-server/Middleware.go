package main

import (
	"github.com/gin-gonic/gin"
)

func AccessCROSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		// 处理js-ajax跨域问题
		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		// w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		c.Next()
	}
}
