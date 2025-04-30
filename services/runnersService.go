package services

import (
	"gopr/models"
	"gopr/repositories"
	_ "log"

	// "log"
	"net/http"
	"strconv"
	"time"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)
	if responseErr != nil {
		// log.Println("Error")
		return nil, responseErr
	}
	return rs.runnersRepository.CreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	responseErr := validateRunnerId(runner.ID)
	if responseErr != nil {
		return responseErr
	}
	responseErr = validateRunner(runner)
	if responseErr != nil {
		return responseErr
	}
	return rs.runnersRepository.UpdateRunner(runner)
}

func (rs RunnersService) DeleteRunner(runnerID string) *models.ResponseError {
	responseErr := validateRunnerId(runnerID)
	if responseErr != nil {
		return responseErr
	}
	return rs.runnersRepository.DeleteRunner(runnerID)
}

func (rs RunnersService) GetRunner(runnerID string) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunnerId(runnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	runner, responseErr := rs.runnersRepository.GetRunner(runnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	results, responseErr := rs.resultsRepository.GetAllRunnersResults(runner.ID)
	if responseErr != nil {
		return nil, responseErr
	}
	runner.Results = results
	return runner, nil

}

func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
	if country == "" && year == "" {
		return nil, &models.ResponseError{
			Message: "No year and country",
			Status:  http.StatusBadRequest,
		}
	}
	if country != "" {
		return rs.runnersRepository.GetRunnersByCountry(country)
	}
	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid integer",
				Status:  http.StatusBadRequest,
			}
		}

		if intYear < 0 || intYear > time.Now().Year() {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		return rs.runnersRepository.GetRunnersByYear(intYear)
	}
	return rs.runnersRepository.GetAllRunners()
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Age < 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func validateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}
