package service

import (
	"time"
)

func (s *Service) Activate() error {
	s.isActive = true
	return nil
}

func (s *Service) Deactivate() error {
	s.isActive = false
	return nil
}

func (s *Service) Rename(newName Name) error {
	if valid, err := newName.IsValid(); !valid {
		return err
	}

	s.name = newName
	s.updatedAt = time.Now().UTC()
	return nil
}

func (s *Service) UpdateDescription(newDescription Description) error {
	if valid, err := newDescription.IsValid(); !valid {
		return err
	}

	s.description = newDescription
	s.updatedAt = time.Now().UTC()
	return nil
}
