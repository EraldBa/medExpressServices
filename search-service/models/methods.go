package models

import "time"

func (s *SearchEntry) AddDefaultData() {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (p *PDFEntry) AddDefaultData() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}
