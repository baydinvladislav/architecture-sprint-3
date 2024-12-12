package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"user-service/schemas/web"
	"user-service/shared"
)

func TestVerifyConnectionService(t *testing.T) {
	appConfig := shared.NewAppSettings()

	// init verify service
	verifyService := NewVerifyConnectionService(appConfig.MinHomeSquare, appConfig.MaxHomeSquare)

	// prepare input data with different test cases
	validUser := &web.UserOut{
		ID:       uuid.New(),
		Username: "testuser",
	}

	invalidUser := &web.UserOut{
		ID:       uuid.New(),
		Username: "", // empty username makes verification is failed
	}

	validHouse := &web.HouseOut{
		ID:      uuid.New(),
		Address: "Valid Address",
		Square:  50.0,
		UserID:  uuid.New(),
	}

	tooSmallHouse := &web.HouseOut{
		ID:      uuid.New(),
		Address: "Too Small House",
		Square:  10.0, // less than min square
		UserID:  uuid.New(),
	}

	tooLargeHouse := &web.HouseOut{
		ID:      uuid.New(),
		Address: "Too Large House",
		Square:  300.0, // greater than max square
		UserID:  uuid.New(),
	}

	// test table
	tests := []struct {
		name           string
		user           *web.UserOut
		house          *web.HouseOut
		expectedResult bool
	}{
		{
			name:           "Valid user and house",
			user:           validUser,
			house:          validHouse,
			expectedResult: true,
		},
		{
			name:           "Invalid user",
			user:           invalidUser,
			house:          validHouse,
			expectedResult: false,
		},
		{
			name:           "House too small",
			user:           validUser,
			house:          tooSmallHouse,
			expectedResult: false,
		},
		{
			name:           "House too large",
			user:           validUser,
			house:          tooLargeHouse,
			expectedResult: false,
		},
		{
			name:           "Invalid user and house",
			user:           invalidUser,
			house:          tooSmallHouse,
			expectedResult: false,
		},
	}

	// tests execution
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := verifyService.VerifyModuleConnection(test.user, test.house)
			require.NoError(t, err)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
