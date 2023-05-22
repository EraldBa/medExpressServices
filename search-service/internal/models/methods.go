package models

import "time"

func (s *Times) AddDefaultData() {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}
