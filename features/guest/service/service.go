package service

import (
	"greenenvironment/features/guest"
)

type GuestService struct {
	guestRepo guest.GuestRepostoryInterface
}

func NewGuestService(guestRepo guest.GuestRepostoryInterface) guest.GuestServiceInterface {
	return &GuestService{guestRepo: guestRepo}
}

func (gs *GuestService) GetGuestProduct() (guest.Guest, error) {
	return gs.guestRepo.GetGuests()
}
