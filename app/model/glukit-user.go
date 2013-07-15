package model

import (
	"lib/goauth2/oauth"
	"time"
)

// Represents a GlukitUser profile
type GlukitUser struct {
	Email          string      `datastore:"email"`
	FirstName      string      `datastore:"firstName"`
	LastName       string      `datastore:"lastName"`
	DateOfBirth    time.Time   `datastore:"birthdate"`
	DiabetesType   string      `datastore:"diabetesType"`
	Timezone       string      `datastore:"timezoneId"`
	LastUpdated    time.Time   `datastore:"lastUpdated"`
	MostRecentRead time.Time   `datastore:"mostRecentRead"`
	Token          oauth.Token `datastore:"token",noindex`
	// This should always remain there, if we lose it, we have to force to get a brand-new one
	RefreshToken string      `datastore:"refreshToken",noindex`
	Score        GlukitScore `datastore:"score"`
}

// Represents a GlukitScore value
type GlukitScore struct {
	Value      int64     `datastore:"value"`
	UpperBound time.Time `datastore:"upperBound"`
	UpdatedAt  time.Time `datastore:"updatedAt"`
}
