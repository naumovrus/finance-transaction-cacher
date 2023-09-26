package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naumovrus/finance-transaction-api/internal/entity"
	"github.com/sirupsen/logrus"
)

func (h *Handler) createWallet(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	id, err := h.services.Money.CreateWallet(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type amountRequest struct {
	Amount string `json:"amount"`
}

func (h *Handler) topUp(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input amountRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(input.Amount, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var id int

	id, err = h.services.Money.TopUp(userId, amount)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id_transaction": id,
	})
}

func (h *Handler) takeOut(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input amountRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(input.Amount, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var id int
	id, err = h.services.Money.TakeOut(userId, amount)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id_transaction": id,
	})
}

type sendMoneyRequest struct {
	UserIdTo string `json:"user_id_to"`
	Amount   string `json:"amount"`
}

func (h *Handler) send(c *gin.Context) {
	userIdFrom, err := getUserId(c)
	if err != nil {
		return
	}

	// jstream := `{"user_id_to": "1", "amount": "12.00"}`

	var input sendMoneyRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// err = json.NewDecoder(strings.NewReader(jstream)).Decode(&input)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	amountstr := input.Amount

	amount, err := strconv.ParseFloat(amountstr, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var lastId int
	lastId, err = h.services.Money.GetLastTransactionSend()
	log.Printf("last_id: %v", lastId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var cache entity.TransactionSend
	cache.Id = lastId + 1
	cache.UserIdFrom = userIdFrom
	userIdTo, err := strconv.Atoi(input.UserIdTo)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	cache.UserIdTo = userIdTo
	cache.Time = time.Now()
	h.redisCache.SetTS(cache.Time.String(), cache)
	var id int
	// id, err = h.services.Money.Send(userIdFrom, userIdTo, amount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Printf("amount send: %v, user_id_from: %v, user_id_to: %v", amount, userIdFrom, userIdTo)
	c.JSON(http.StatusOK, map[string]interface{}{
		"transaction_id": id,
	})
}
