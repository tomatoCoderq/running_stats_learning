package services

import (
	"github.com/stretchr/testify/assert"
	"gopr/models"
	"net/http"
	"testing"
)

func TestValidateRunnerInvalidFirstName(t *testing.T) {
	runner := &models.Runner{
		LastName: "smith",
		Age:      40,
		Country:  "Country",
	}
	responseErr := validateRunner(runner)
	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid first name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunner(t *testing.T) {
	tests := []struct {
		name string
		runner *models.Runner
		want *models.ResponseError
	}{
		{
			name: "Invalid_First_Name",
			runner: &models.Runner{LastName: "smith",
			Age:      40,
			Country:  "Country",
			},
			want: &models.ResponseError {
				Message: "Invalid first name",
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Last_Name",
			runner: &models.Runner{
				FirstName: "John",
				Age:      40,
				Country:  "Country",
			},
			want: &models.ResponseError {
				Message: "Invalid last name",
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Age",
			runner: &models.Runner{
				FirstName: "smith",
				LastName: "smith",
				Age:      -1,
				Country:  "Country",
			},
			want: &models.ResponseError {
				Message: "Invalid age",
				Status: http.StatusBadRequest,
			},
		},
		{
			name: "Valid_Runner",
			runner: &models.Runner{
				FirstName: "smith",
				LastName: "aga",
				Age:      23,
				Country:  "Country",
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func (t *testing.T) {
			responseErr := validateRunner(test.runner)
			assert.Equal(t, test.want, responseErr)
		})
	}
}