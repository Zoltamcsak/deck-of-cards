package handler

import (
	custErr "github.com/deck/internal/app/error"
	"github.com/deck/internal/app/model"
	"github.com/deck/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DeckHandler struct {
	service service.DeckService
}

func NewDeckHandler(service service.DeckService) *DeckHandler {
	return &DeckHandler{service: service}
}

func (h *DeckHandler) CreateDeck(ctx *gin.Context) {
	var err error
	cards := ctx.Query("cards")
	shuffledParam := ctx.Query("shuffled")
	shuffled := false
	if len(shuffledParam) > 0 {
		shuffled, err = strconv.ParseBool(shuffledParam)
		if err != nil {
			serveHttpError(ctx, custErr.New(http.StatusBadRequest, "shuffled must be boolean"))
			return
		}
	}

	req := model.CreateDeckRequest{
		Shuffled: shuffled,
		Cards:    cards,
	}
	deck, err := h.service.CreateDeck(req)
	if err != nil {
		serveHttpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, deck)
}

func (h *DeckHandler) GetDeckById(ctx *gin.Context) {
	id := ctx.Param("id")
	deck, err := h.service.GetDeckById(id)
	if err != nil {
		serveHttpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, deck)
}

func (h *DeckHandler) DrawCards(ctx *gin.Context) {
	id := ctx.Param("id")
	countParam := ctx.Query("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		serveHttpError(ctx, custErr.New(http.StatusBadRequest, "count must be a number"))
		return
	}
	cards, err := h.service.DrawCards(id, count)
	if err != nil {
		serveHttpError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, cards)
}

func (h *DeckHandler) InitRoutes(engine *gin.Engine) {
	engine.POST("/decks", h.CreateDeck)
	engine.GET("/decks/:id", h.GetDeckById)
	engine.PUT("/decks/:id/cards", h.DrawCards)
}
