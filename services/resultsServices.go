package services

import (
	"gopr/models"
	"gopr/repositories"
	_"log"
	"net/http"
	_"strconv"
	"time"
)

type ResultsService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewResultsService(runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *ResultsService {
	return &ResultsService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid RunnerID",
			Status:  http.StatusBadRequest,
		}
	}
	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid Race Results",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid Location",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid Position",
			Status:  http.StatusBadRequest,
		}
	}
	if result.ID == " " {
		return nil, &models.ResponseError{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Year < 0 || result.Year > time.Now().Year() {
		return nil, &models.ResponseError{
			Message: "Invalid Year",
			Status:  http.StatusBadRequest,
		}
	}

	raceResult, err := ParseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid Race Result",
			Status:  http.StatusBadRequest,
		}
	}

	response, responseError := rs.resultsRepository.CreateResult(result)
	if responseError != nil {
		return nil, responseError
	}

	runner, responseError := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseError != nil {
		return nil, responseError
	}
	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	//update the best runner's result
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := ParseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse",
				Status:  http.StatusInternalServerError,
			}
		}
		if personalBest > raceResult {
			runner.PersonalBest = result.RaceResult
		}
	}

	//update the best runner's seasonal result
	if runner.SeasonBest == "" {
		runner.SeasonBest = result.RaceResult
	} else {
		seasonBest, err := ParseRaceResult(runner.SeasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse",
				Status:  http.StatusInternalServerError,
			}
		}
		if raceResult < seasonBest {
			runner.SeasonBest = raceResult.String()
		}
	}

	responseError = rs.runnersRepository.UpdateRunnerResults(runner)
	if err != nil {
		return nil, responseError
	}
	return response, nil
}

func (rs ResultsService) DeleteResult(resultID string) *models.ResponseError {
	if resultID == "" {
		return &models.ResponseError{
			Message: "invalid resultID",
			Status:  http.StatusBadRequest,
		}
	}

	err := repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
	if err != nil {
		return &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}

	result, responseErr := rs.resultsRepository.DeleteResult(resultID)
	if responseErr != nil {
		return responseErr
	}

	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		repositories.RollBackTransaction(rs.runnersRepository, rs.resultsRepository)
		return responseErr
	}

	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			repositories.RollBackTransaction(rs.runnersRepository, rs.resultsRepository)
			return responseErr
		}
		runner.PersonalBest = personalBest
	}

	if runner.SeasonBest == result.RaceResult && result.Year == time.Now().Year() {
		seasonBest, responseErr := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, time.Now().Year())
		if responseErr != nil {
			repositories.RollBackTransaction(rs.runnersRepository, rs.resultsRepository)
			return responseErr
		}
		runner.SeasonBest = seasonBest
	}

	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		repositories.RollBackTransaction(rs.runnersRepository, rs.resultsRepository)
		return responseErr
	}
	repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)
	return nil

}

func ParseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(timeString[0:2] + "h" + timeString[3:5] + "m" + timeString[6:8] + "s")
}

// func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError){
// 	responseErr := validateRunner(runner)
// 	if responseErr != nil {
// 		return nil, responseErr
// 	}
// 	return rs.runnersRepository.CreateRunner(runner)
// }

// func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError{
// 	responseErr := validateRunnerId(runner.ID)
// 	if responseErr != nil {
// 		return responseErr
// 	}
// 	responseErr = validateRunner(runner)
// 	if responseErr != nil {
// 		return responseErr
// 	}
// 	return rs.runnersRepository.UpdateRunner()
// }

// func (rs RunnersService) DeleteRunner(runnerID string) *models.ResponseError{
// 	responseErr := validateRunnerId(runnerID)
// 	if responseErr != nil {
// 		return responseErr
// 	}
// 	return rs.runnersRepository.DeleteRunner(runnerID)
// }

// func (rs RunnersService) GetRunner(runnerID string) (*models.Runner, *models.ResponseError){
// 	responseErr := validateRunnerId(runnerID)
// 	if responseErr != nil {
// 		return nil, responseErr
// 	}
// 	runner, responseErr := rs.runnersRepository.GetRunner(runnerID)
// 	if responseErr != nil {
// 		return nil, responseErr
// 	}
// 	results, responseErr := rs.resultsRepository.GetAllRunnersResults(runner)
// 	if responseErr != nil {
// 		return nil, responseErr
// 	}
// 	runner.Results = results
// 	return runner, nil

// }

// func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError){
// 	if country == "" && year == "" {
// 		return nil, &models.ResponseError {
// 			Message: "No year and country",
// 			Status: http.StatusBadRequest,
// 		}
// 	}
// 	if country != "" {
// 		return rs.runnersRepository.GetRunnerByCountry(country)
// 	}
// 	if year != "" {
// 		intYear, err := strconv.Atoi(year)
// 		if err != nil {
// 			return nil, &models.ResponseError {
// 				Message: "Invalid integer",
// 				Status: http.StatusBadRequest,
// 			}
// 		}

// 		if intYear < 0 || intYear > time.Now().Year() {
// 			return nil, &models.ResponseError {
// 				Message: "Invalid year",
// 				Status: http.StatusBadRequest,
// 			}
// 		}

// 		return rs.runnersRepository.GetRunnerByYear(intYear)
// 	}
// 	return rs.runnersRepository.GetAllRunners()
// }

// func validateRunner(runner *models.Runner) *models.ResponseError {
// 	if runner.FirstName == "" {
// 		return &models.ResponseError {
// 			Message: "Invalid first name",
// 			Status: http.StatusBadRequest,
// 		}
// 	}

// 	if runner.LastName == "" {
// 		return &models.ResponseError {
// 			Message: "Invalid last name",
// 			Status: http.StatusBadRequest,
// 		}
// 	}

// 	if runner.Age < 16 || runner.Age > 125 {
// 		return &models.ResponseError {
// 			Message: "Invalid age",
// 			Status: http.StatusBadRequest,
// 		}
// 	}
// 	return nil;
// }

// func validateRunnerId(runnerId string) *models.ResponseError {
// 	if runnerId == "" {
// 		return &models.ResponseError {
// 			Message: "Invalid ID",
// 			Status: http.StatusBadRequest,
// 		}
// 	}
// 	return nil
// }
