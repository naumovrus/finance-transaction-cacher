package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	err = h.services.Money.TopUp(userId, amount)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Updated",
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

	err = h.services.Money.TakeOut(userId, amount)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Updated",
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

	var input sendMoneyRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// var idtx int

	amountstr := input.Amount

	amount, err := strconv.ParseFloat(amountstr, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userIdToStr := input.UserIdTo
	userIdTo, err := strconv.Atoi(userIdToStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Money.Send(userIdFrom, userIdTo, amount)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Printf("amount send: %v, user_id_from: %v, user_id_to: %v", amount, userIdFrom, userIdTo)
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"amount": amount,
	})
}
