package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopr/models"
	"gopr/services"
	"io"
	"log"
	"net/http"	
	_"github.com/google/uuid"

)

type RunnersController struct {
	runnersService *services.RunnersService
}

func NewRunnersController(runnersService *services.RunnersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
	}
}

func (rh RunnersController) CreateRunner(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while reading", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := rh.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		log.Println(responseErr)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)

}

func (rh RunnersController) UpdateRunner(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading "+"update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling "+"update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	responseErr := rh.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)

}

func (rh RunnersController) DeleteRunner(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	responseErr := rh.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rh RunnersController) GetRunner(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	log.Print(runnerId)
	response, responseErr := rh.runnersService.GetRunner(runnerId)
	if responseErr != nil {
		log.Printf("Error while getting runner %v", responseErr)
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rh RunnersController) GetRunnersBatch(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")
	response, responseErr := rh.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return 
	}
	ctx.JSON(http.StatusOK, response)
}
