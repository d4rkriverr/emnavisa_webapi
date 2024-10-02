package calllog

import (
	"database/sql"
	"emnavisa/webserver/infrastructure/kernel"
	"emnavisa/webserver/infrastructure/models"
	"fmt"
	"log"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(app *kernel.Application) *Service {
	return &Service{
		db: app.Database,
	}
}

func (s *Service) GetAllCallsByAgent(agentName string, date time.Time) ([]models.CallLog, error) {
	calls := []models.CallLog{}

	// SQL query to select all calls where the agent matches
	query := `SELECT id, cin, name, phone, requested_job, requested_country, created_at, platform, agent, call_status, notes 
	          FROM call_logs WHERE agent = $1  AND DATE(created_at) = $2 ORDER BY created_at DESC`

	// Execute the query
	rows, err := s.db.Query(query, agentName, date)
	if err != nil {
		log.Println("Error executing query:", err)
		return calls, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the CallLog slice
	for rows.Next() {
		var call models.CallLog
		err := rows.Scan(
			&call.ID,
			&call.CIN,
			&call.Name,
			&call.Phone,
			&call.RequestedJob,
			&call.RequestedCountry,
			&call.CreatedAt,
			&call.Platform,
			&call.Agent,
			&call.CallStatus,
			&call.Notes,
		)
		if err != nil {
			return calls, err
		}
		calls = append(calls, call)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		return calls, err
	}

	return calls, nil
}

func (s *Service) CreateNewCallLog(callLog models.CallLog) error {
	// SQL query to insert a new call log
	query := `
		INSERT INTO call_logs 
		(cin, name, phone, requested_job, requested_country, created_at, platform, agent, call_status, notes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// Execute the query
	_, err := s.db.Exec(query,
		callLog.CIN,
		callLog.Name,
		callLog.Phone,
		callLog.RequestedJob,
		callLog.RequestedCountry,
		callLog.CreatedAt,
		callLog.Platform,
		callLog.Agent,
		callLog.CallStatus,
		callLog.Notes)

	if err != nil {
		log.Println("Error inserting new call log:", err)
		return err
	}

	return nil
}

func (s *Service) UpdateCallLog(agent string, updatedCallLog models.CallLog) error {
	query := `
        UPDATE call_logs
        SET cin = $1, name = $2, phone = $3, requested_job = $4, 
            requested_country = $5, call_status = $6, platform = $7, 
            notes = $8
        WHERE agent = $9 AND id = $10
    `
	_, err := s.db.Exec(query,
		updatedCallLog.CIN,
		updatedCallLog.Name,
		updatedCallLog.Phone,
		updatedCallLog.RequestedJob,
		updatedCallLog.RequestedCountry,
		updatedCallLog.CallStatus,
		updatedCallLog.Platform,
		updatedCallLog.Notes,
		agent,
		updatedCallLog.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update call log: %w", err)
	}

	return nil
}
