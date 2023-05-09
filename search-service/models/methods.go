package models

import "time"

func (s *SearchEntry) AddDefaultData() {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *PDFEntry) AddDefaultData() {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}
