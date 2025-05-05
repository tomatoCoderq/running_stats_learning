package controllers

import (
	"database/sql"
	"encoding/json"
	"gopr/models"
	"gopr/repositories"
	"gopr/services"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initTestRouter(dbHandler *sql.DB) *gin.Engine{ 
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	runnersController := NewRunnersController(runnersService)
	router := gin.Default()

	router.GET("/runner", runnersController.GetRunnersBatch)
	return router
}

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()
	columns := []string{"id", "first_name", "last_name", "age", "is_acitve", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).AddRow(
			"1", "John", "Smith", 30, true, "US", "02:00:30", "2:13:03").AddRow(
				"1", "Brzno", "Brznovich", 23, true, "Serbia", "02:00:30", "2:13:03"))
	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	recorde := httptest.NewRecorder()
	router.ServeHTTP(recorde, request)
	assert.Equal(t, http.StatusOK, recorde.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorde.Body.Bytes(), &runners)
	assert.NotEmpty(t, runners)
	assert.Equal(t, 2, len(runners))
}

