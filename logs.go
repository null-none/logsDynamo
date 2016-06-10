package main

import (
    "github.com/gin-gonic/gin"
    "gopkg.in/go-playground/validator.v8"
    "github.com/goamz/goamz/aws"
    "github.com/goamz/goamz/dynamodb"
    "time"
    "log"
)

type Logs struct {
    UnixTimestamp   string     `validate:"required"`
    TimeZone        string     `validate:"required"`
    URL             string     `validate:"required"`
    Browser         string     `validate:"required"`
    BrowserVersion  string     `validate:"required"`
    IP              string     `validate:"required"`
    OS              string     `validate:"required"`
    Display         string     `validate:"required"`
    Flash           string     `validate:"required"`
    Device          string     `validate:"required"`
    JavaScript      string     `validate:"required"`
    UserAgent       string     `validate:"required"`
    Lang            string     `validate:"required"`
}

var validate *validator.Validate

func main() {

    r := gin.Default()
    r.POST("/logs", func(c *gin.Context) {
        auth, err := aws.EnvAuth()
        if err != nil {
            log.Panic(err)
        }

        config := &validator.Config{TagName: "validate"}
        validate = validator.New(config)

        log := &Logs{
          UnixTimestamp:     c.PostForm("unixtimestamp"),
          TimeZone:          c.PostForm("time_zone"),
          URL:               c.PostForm("url"),
          Browser:           c.PostForm("browser"),
          IP:                c.DefaultPostForm("ip", "localhost"),
          UserAgent:         c.DefaultPostForm("user_agent", "user_agent"),
          BrowserVersion:    c.PostForm("browser_version"),
          OS:                c.PostForm("os"),
          Display:           c.PostForm("display"),
          Flash:             c.PostForm("flash"),
          Device:            c.PostForm("device"),
          JavaScript:        c.PostForm("javascript"),
          Lang:              c.PostForm("lang"),
        }

        errs := validate.Struct(log)

        if errs != nil {
          c.JSON(400, gin.H{
              "message": errs,
              "status": "error",
          })
        } else {
          ddbs := dynamodb.Server{auth, aws.EUCentral}
          pkattr := dynamodb.NewStringAttribute("id", "")
        	pk := dynamodb.PrimaryKey{pkattr, nil}
        	table := dynamodb.Table{&ddbs, "logs", pk}
          ats, err := dynamodb.MarshalAttributes(visit)
          if err != nil {
              log.Panic(err)
          }
          id := int32(time.Now().Unix())
          row , err := table.PutItem(string(id), "", ats)
          if err != nil {
              log.Panic(err)
          }
          c.JSON(200, gin.H{
              "status": "ok",
          })
        }
    })
    r.Run() // listen and server on 0.0.0.0:8080
}
