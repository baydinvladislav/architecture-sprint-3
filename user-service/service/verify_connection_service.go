package service

import (
	"user-service/schemas/web"
)

type VerifyConnectionService struct {
	minSquare float64
	maxSquare float64
}

func NewVerifyConnectionService(minSquare float64, maxSquare float64) *VerifyConnectionService {
	return &VerifyConnectionService{
		minSquare: minSquare,
		maxSquare: maxSquare,
	}
}

func (r *VerifyConnectionService) VerifyModuleConnection(
	verifyingUser *web.UserOut,
	verifyingHouse *web.HouseOut,
) (bool, error) {
	if r.verifyUser(verifyingUser) == true && r.verifyHouse(verifyingHouse) == true {
		return true, nil
	}
	return false, nil
}

func (r *VerifyConnectionService) verifyUser(verifyingUser *web.UserOut) bool {
	if verifyingUser.Username == "" {
		return false
	}
	return true
}

func (r *VerifyConnectionService) verifyHouse(verifyingHouse *web.HouseOut) bool {
	if verifyingHouse.Square < r.minSquare || verifyingHouse.Square > r.maxSquare {
		return false
	}
	return true
}
