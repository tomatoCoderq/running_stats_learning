package repositories

import (
	"database/sql"
	"gopr/models"
	"net/http"
	_"strconv"
)

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB) *RunnersRepository {
	return &RunnersRepository{
		dbHandler: dbHandler,
	}
}

func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	query := `INSERT into runners(first_name, last_name, age, country) values ($1, $2, $3, $4) returning id`
	rows, err := rr.dbHandler.Query(query, runner.FirstName, runner.LastName, runner.Age, runner.Country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	

	defer rows.Close()
	var runnerId string
	for rows.Next() {
		err := rows.Scan(&runnerId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:        runnerId,
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  true,
		Country:   runner.Country,
	}, nil

}

func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	query := `UPDATE runners SET 
	first_name=$1,
	second_name=$2, 
	age=$3,
	country=$4,
	WHERE id=$5`

	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country, runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "No user found",
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (rr RunnersRepository) UpdateRunnerResults(runner *models.Runner) *models.ResponseError {
	return nil
}

func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {
	query := `UPDATE runners
	SET is_active='false'
	WHERE id=$1`

	res, err := rr.dbHandler.Exec(query, runnerId)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "No user found",
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	query := `SELECT * FROM runners where id=$1`
	// var runnerIdUU uint;
	// runnerIdUU, _ := strconv.Atoi(runnerId)


	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool
	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		IsActive:     isActive,
		Country:      country,
		PersonalBest: personalBest.String,
		SeasonBest:   seasonBest.String,
	}, nil
}

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {
	query := `SELECT * FROM runners`
	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runners = append(runners, &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		})
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil

}

func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	query := `SELECT * FROM runners where country=$1 AND is_active='true' order by personal_best LIMIT 10`
	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runners = append(runners, &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		})
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) GetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {
	query := `SELECT runners.id, runners.first_name, runners.last_name, runners.age, runners.is_alive, 
	runners.personal_best, results.race_result
	FROM runners INNER JOIN 
	(
		SELECT runner_id,
		MIN(race_result) as race_result
		FROM results
		WHERE year=$1
		GROUPED BY runner_id
	)
	results ON runners.id = results.runner_id
	ORDER BY results.race_result
	LIMIT 10`

	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil ,&models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runners = append(runners, &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		})
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}
