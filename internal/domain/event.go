package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Title       string 
	Description string   
	Location    string  
	StartTime   time.Time 
	EndTime     *time.Time 
	CreatedAt   time.Time 
}

type Seat struct {
	ID      uuid.UUID 
	EventID uuid.UUID 
	Row     int      
	Number  int       
	Price   int       
	Booking *Booking  
}
