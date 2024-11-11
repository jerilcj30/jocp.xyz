package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type LeadsTable struct {
	ID   int
	Lead json.RawMessage
}

type LeadsService struct {
	DB *sql.DB
}

func (ls *LeadsService) Create(data json.RawMessage) (*LeadsTable, error) {

	lead := LeadsTable{
		Lead: data,
	}

	row := ls.DB.QueryRow(`INSERT INTO leads (lead) VALUES ($1) RETURNING id`, data)
	err := row.Scan(&lead.ID)
	if err != nil {
		return nil, fmt.Errorf("Create Lead:%w", err)
	}

	return &lead, nil
}
