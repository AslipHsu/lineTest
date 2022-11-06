package comon

import (
	dbM "docker_test/database/mgo"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ReceiveLineMessage struct {
	Destination string        `json:"destination"`
	Events      []interface{} `json:"events"`
}

type LineSendReq struct {
	ResText []string `json:"resText"`
	SendTo  string   `json:"sendTo"`
}

func LineReceive(c *gin.Context) {
	fmt.Println("line message received")
	appG := Gin{C: c}
	input := &ReceiveLineMessage{}
	err := c.Bind(input)
	if err != nil {
		appG.response(-1, 400, "parse err", nil)
		fmt.Println("parse err:", err)
		return
	}

	//Parse input
	fmt.Println("Events:  ", input.Events)
	event := input.Events[0].(map[string]interface{})
	mes := event["message"]
	source := event["source"]
	// fmt.Println("mes:  ", mes)
	// fmt.Println("source:  ", source)
	userid := source.(map[string]interface{})["userId"].(string)
	messageData := &LineMessageData{
		UserId:    userid,
		MessageId: mes.(map[string]interface{})["id"].(string),
		Text:      mes.(map[string]interface{})["text"].(string),
		Timestamp: event["timestamp"].(float64),
	}

	//reply
	replyToken := event["replyToken"].(string)
	fmt.Println("replyToken:  ", replyToken)
	resText, _ := makeReplyText("received message!!", "")
	err = replyText(replyToken, resText)
	if err != nil {
		appG.response(-1, 400, "replyText err", nil)
		fmt.Println("replyText err:", err)
		return
	}

	//save db
	err = addUserData(userid)
	if err != nil {
		appG.response(-1, 400, "addUserData err", nil)
		fmt.Println("addUserData err:", err)
		return
	}
	err = addMessageData(messageData)
	if err != nil {
		appG.response(-1, 400, "addMessageData err", nil)
		fmt.Println("addMessageData err:", err)
		return
	}

	appG.response(0, 200, "success", nil)
}

func LineSend(c *gin.Context) {
	fmt.Println("send line message !!")
	appG := Gin{C: c}
	input := &LineSendReq{}
	err := c.Bind(input)
	if err != nil {
		appG.response(-1, 400, "parse err", nil)
		fmt.Println("parse err:", err)
		return
	}

	resText := make([]interface{}, 0)
	for _, i := range input.ResText {
		if len(i) > 0 {
			resText, _ = makeReplyText(resText, i)
		}
	}
	// sendto := ""

	err = pushText(input.SendTo, resText)
	if err != nil {
		appG.response(-1, 400, "pushText err", nil)
		fmt.Println("pushText err:", err)
		return
	}

	appG.response(0, 200, "success", nil)
}

func GetLineMessages(c *gin.Context) {
	fmt.Println("Get Line Messages")
	appG := Gin{C: c}
	id := c.Param("userId")
	page, _ := strconv.Atoi(c.Param("page"))
	limit, _ := strconv.Atoi(c.Param("limit"))

	filter := bson.M{}
	if len(id) > 0 {
		filter["userId"] = id
	}
	input := &List{
		Page:   int64(page),
		Limit:  int64(limit),
		Filter: filter,
		// SortField: "",
		// SortType:  -1,
	}
	data := []LineMessageData{}
	cmd := dbM.SearchCMD{
		DBName: db_name,
		CName:  db_table_message,
	}
	count, err := findAllFromDB(cmd, input, &data)
	if data == nil || err != nil || len(data) <= 0 {
		r := make([]LineMessageData, 0)
		fmt.Println("no data(list)")
		fmt.Println("db err: :", err)
		appG.listResponse(0, 200, "find no data", 0, r)
		return
	}
	//return
	appG.listResponse(0, 200, "success", count, data)
}
