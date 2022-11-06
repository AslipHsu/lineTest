package comon

import (
	"bytes"
	dbM "docker_test/database/mgo"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db_name          = "line"
	db_table_user    = "user"
	db_table_message = "message"

	channelSecrect     = "ad4c2bf6b628a05e585b3f474ce22734"
	channelAccessToken = "VhT40k79sF+H9ahtBqrGlCZ/T2xPpwGveY2jgm7nkG0iPNxbRyUaM+3P765/dspDbthQ02VczfCDzPzXTTWGega3rUL7dZuUs2LbiPAnasOMSZ/BIQHfQib/Zr2JPtMIVBCmrao+hMn7ZNV1frrR/AdB04t89/1O/w1cDnyilFU="
	replyUrl           = "https://api.line.me/v2/bot/message/reply"
	pushUrl            = "https://api.line.me/v2/bot/message/push"
	profileUrl         = "https://api.line.me/v2/bot/profile/{userId}"
)

/////////////////////////////////////////////////////////////////////////

type ReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LineUserData struct {
	UserId      string `json:"userId" bson:"userId"`
	DisplayName string `json:"displayName" bson:"displayName"`
	PictureUrl  string `json:"pictureUrl" bson:"pictureUrl"`
	Language    string `json:"language" bson:"language"`
}

type LineMessageData struct {
	UserId    string  `json:"userId" bson:"userId"`
	MessageId string  `json:"messageId" bson:"messageId"`
	Text      string  `json:"text" bson:"text"`
	Timestamp float64 `json:"timestamp" bson:"timestamp"`
}

// 通用列表查詢格式
type List struct {
	Page      int64  `json:"page"`  //第幾頁
	Limit     int64  `json:"limit"` //一頁多少筆
	Filter    bson.M `json:"filter"`
	SortField string `json:"sort_field"` //要排序的欄位名稱
	SortType  int    `json:"sort_type"`  //排序類型(1正序,-1反序)
}

type Gin struct {
	C *gin.Context
}

// initText可為string or []interface{}
// addText 可為空
// 一則訊息以上用addText增加
// Max: 5
func makeReplyText(initText interface{}, addText string) ([]interface{}, error) {
	resText := make([]interface{}, 1)
	switch t := initText.(type) {
	case string:
		resText[0] = &ReplyMessage{
			Type: "text",
			Text: t,
		}
		if len(addText) > 0 {
			resText2 := &ReplyMessage{
				Type: "text",
				Text: addText,
			}
			resText = append(resText, resText2)
		}

	case []interface{}:
		if len(addText) > 0 {
			resText2 := &ReplyMessage{
				Type: "text",
				Text: addText,
			}
			resText = append(t, resText2)
		}
	default:
		return nil, errors.New("type err")
	}
	return resText, nil
}

func replyText(replyToken string, resText []interface{}) error {
	resT := make(map[string]interface{})
	resT["replyToken"] = replyToken
	resT["messages"] = resText
	jsonStr, _ := json.Marshal(resT)
	client := &http.Client{}
	r, err := http.NewRequest("POST", replyUrl, bytes.NewBuffer(jsonStr)) // URL-encoded payload
	if err != nil {
		fmt.Println("replyText err 1 : ", err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+channelAccessToken)

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("replyText err 2 : ", err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("replyText err 3 : ", err)
		return err
	}
	if string(body) != "{}" {
		e := fmt.Sprintln("pushText string(body):", string(body))
		return errors.New(e)

	}
	return nil
}

func pushText(sendto string, resText []interface{}) error {
	resT := make(map[string]interface{})
	resT["to"] = sendto
	resT["messages"] = resText
	jsonStr, _ := json.Marshal(resT)
	client := &http.Client{}
	r, err := http.NewRequest("POST", pushUrl, bytes.NewBuffer(jsonStr)) // URL-encoded payload
	if err != nil {
		fmt.Println("pushText err 1 : ", err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+channelAccessToken)

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("pushText err 2 : ", err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("pushText err 3 : ", err)
		return err
	}
	if string(body) != "{}" {
		e := fmt.Sprintln("pushText string(body):", string(body))
		return errors.New(e)

	}
	return nil
}

func queryLineUserData(userid string) (*LineUserData, error) {
	client := &http.Client{}
	url := strings.Replace(profileUrl, "{userId}", userid, 1)
	r, err := http.NewRequest("GET", url, nil) // URL-encoded payload
	if err != nil {
		fmt.Println("GET err 1 : ", err)
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+channelAccessToken)

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("GET err 2 : ", err)
		return nil, err
	}
	defer res.Body.Close()
	arrayByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("GET err 3 : ", err)
		return nil, err
	}
	rst := &LineUserData{}
	err = json.Unmarshal(arrayByte, rst)
	if err != nil {
		fmt.Println("GET err 4 : ", err)
		return nil, err
	}
	return rst, nil
}

func addUserData(userId string) error {
	//判斷user是否存在
	dataCount, err := dbM.MgoDB.Count(dbM.SearchCMD{
		DBName: db_name,
		CName:  db_table_user,
		Query:  bson.M{"userId": userId},
	})
	if err != nil {
		return err
	}
	if dataCount <= 0 {
		fmt.Println("save new user")
		//新user存db
		userData, err := queryLineUserData(userId)
		if err != nil {
			return err
		}
		cmd := dbM.SearchCMD{
			DBName: db_name,
			CName:  db_table_user,
			Insert: userData,
		}
		_, err = dbM.MgoDB.Insert(cmd)
		if err != nil {
			return err
		}
	}
	return err
}

func addMessageData(input *LineMessageData) error {
	cmd := dbM.SearchCMD{
		DBName: db_name,
		CName:  db_table_message,
		Insert: input,
	}
	_, err := dbM.MgoDB.Insert(cmd)
	if err != nil {
		return err
	}
	return err
}

// 通用列表查詢,包含模糊查詢
// cmd帶DBName,DBName
// batch 代入自訂回傳結構,ex: batch:=&[]UserData{}
// filter有值就模糊查詢,沒有就查全部
func findAllFromDB(cmd dbM.SearchCMD, input *List, batch interface{}) (int64, error) {
	var skip int64
	var err error
	if input.Page != int64(0) {
		skip = (input.Page - 1) * input.Limit
	} else {
		skip = 0
	}

	findOptions := options.Find()
	input.Filter["is_delete"] = bson.M{"$ne": true} //避免查到已刪除資料
	if len(input.SortField) > 0 {
		findOptions.SetSort(bson.M{input.SortField: input.SortType})
	}
	findOptions.SetLimit(input.Limit)
	findOptions.SetSkip(skip)

	cmd.Query = input.Filter
	count, err := dbM.MgoDB.Count(cmd)
	if err != nil {
		return count, err
	}
	err = dbM.MgoDB.FindOptions(cmd, batch, *findOptions)
	return count, err
}

// 一般類返回
func (g *Gin) response(status int, errCode int, msg string, data interface{}) {
	if data != nil {
		g.C.JSON(http.StatusOK, gin.H{
			"status": status,
			"msg":    msg,
			"data":   data,
		})
	} else {
		g.C.JSON(http.StatusOK, gin.H{
			"status": status,
			"msg":    msg,
		})
	}
	return
}

// 列表類返回
func (g *Gin) listResponse(status int, errCode int, msg string, count int64, data interface{}) {
	if data != nil {
		g.C.JSON(http.StatusOK, gin.H{
			"status": status,
			"msg":    msg,
			"count":  count,
			"data":   data,
		})
	} else {
		g.C.JSON(http.StatusOK, gin.H{
			"status": status,
			"msg":    msg,
			"count":  count,
		})
	}
	return
}
