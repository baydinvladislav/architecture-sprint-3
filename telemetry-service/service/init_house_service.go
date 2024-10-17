package service

import (
	"fmt"
	"log"
	"telemetry-service/repository"
	"telemetry-service/schemas"
)

type InitHouseService struct {
	houseRepository *repository.HouseRepository
}

func NewInitHouseService(houseRepository *repository.HouseRepository) *InitHouseService {
	return &InitHouseService{
		houseRepository: houseRepository,
	}
}

func (s *InitHouseService) ProcessEvent(event schemas.Event) error {
	data, ok := event.Payload.(schemas.InstallModuleToHousePayload)
	if !ok {
		return fmt.Errorf("invalid payload for InstallModuleToHouse event")
	}

	// some logic with data from event ...

	err := s.houseRepository.InsertHouse(event)
	if err != nil {
		return err
	}

	log.Println("msg data: ", data)
	return nil
}
