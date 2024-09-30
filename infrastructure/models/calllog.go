package models

import "time"

type CallLog struct {
	ID               int       `json:"id"`                // Unique identifier for the call log
	CIN              string    `json:"cin"`               // National Identification Number
	Name             string    `json:"name"`              // Name of the client
	Phone            string    `json:"phone"`             // Phone number
	RequestedJob     string    `json:"requested_job"`     // Job requested by the client
	RequestedCountry string    `json:"requested_country"` // Country where the job is requested
	CallStatus       string    `json:"call_status"`       // Call outcome (e.g., reached, not reached)
	Platform         string    `json:"platform"`          // Platform used for the contract (e.g., mobile, web)
	Notes            string    `json:"notes"`             // Additional notes from the call center agent
	Agent            string    `json:"agent"`             // Name of the agent handling the call
	CreatedAt        time.Time `json:"created_at"`        // Date when the call was logged
}
