package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type CheckSpamRequest struct {
	AccessKey int64 `json:"access_key"`
	CheckSpam int   `json:"check_spam"`
}

type Request struct {
	MessageId int64 `json:"message_id"`
}

func (h *Handler) issueCheckSpam(c *gin.Context) {
	var input Request
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Началась обработка сообщения №", input.MessageId)

	c.Status(http.StatusOK)

	go func() {
		time.Sleep(3 * time.Second)
		sendCheckSpamRequest(input)
	}()
}

func sendCheckSpamRequest(request Request) {

	var checkSpam = -1
	if rand.Intn(10)%10 >= 3 {
		checkSpam = rand.Intn(2)

		fmt.Printf("Сообщение №%d обработано\n", request.MessageId)
		fmt.Println("Проверка на спам: ", checkSpam)
	} else {
		fmt.Printf("Сообщение №%d не получилось обработать\n", request.MessageId)
	}

	answer := CheckSpamRequest{
		AccessKey: 123,
		CheckSpam: checkSpam,
	}

	client := &http.Client{}

	jsonAnswer, _ := json.Marshal(answer)
	bodyReader := bytes.NewReader(jsonAnswer)

	requestURL := fmt.Sprintf("http://127.0.0.1:8000/api/messages/%d/update_check_spam/", request.MessageId)

	req, _ := http.NewRequest(http.MethodPut, requestURL, bodyReader)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer response.Body.Close()
}
