package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Workspace struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Floor   string `json:"floor_id"`
	Props   Attrs  `json:"properties"`
	Details string `json:"details"`
}

func (this *Workspace) Equal(other *Workspace) bool {
	return this.ID == other.ID && this.Name == other.Name && this.Floor == other.Floor
}

type Booking struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
	Cancelled   bool      `json:"cancelled"`
	CreatedBy   string    `json:"created_by"`
}

func (this *Booking) Equal(other *Booking) bool {
	return this.ID == other.ID && this.WorkspaceID == other.WorkspaceID &&
		this.UserID == other.UserID && this.StartDate == other.StartDate &&
		this.EndDate == other.EndDate && this.Cancelled == other.Cancelled && this.CreatedBy == other.CreatedBy
}

type ExpandedBooking struct {
	Booking
	WorkspaceName string `json:"workspace_name"`
	UserName      string `json:"user_name"`
	FloorID       string `json:"floor_id"`
	FloorName     string `json:"floor_name"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	IsAdmin    bool   `json:"is_admin"`
	Email      string `json:"email"`
}

func (this *User) Equal(other *User) bool {
	return this.ID == other.ID && this.Name == other.Name &&
		this.Email == other.Email && this.Department == other.Department &&
		this.IsAdmin == other.IsAdmin
}

type Floor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Address     string `json:"address"`
}

func (this *Floor) Equal(other *Floor) bool {
	return this.ID == other.ID && this.Name == other.Name &&
		this.DownloadURL == other.DownloadURL
}

type Offering struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
	Cancelled   bool      `json:"cancelled"`
	CreatedBy   string    `json:"created_by"`
}

func (this *Offering) Equal(other *Offering) bool {
	return this.ID == other.ID && this.WorkspaceID == other.WorkspaceID &&
		this.UserID == other.UserID && this.StartDate.Equal(other.StartDate) &&
		this.EndDate == other.EndDate && this.Cancelled == other.Cancelled && this.CreatedBy == other.CreatedBy
}

type ExpandedOffering struct {
	Offering
	WorkspaceName string `json:"workspace_name"`
	UserName      string `json:"user_name"`
	FloorID       string `json:"floor_id"`
	FloorName     string `json:"floor_name"`
}

func (this *ExpandedOffering) Equal(other *ExpandedOffering) bool {
	return this.ID == other.ID && this.WorkspaceID == other.WorkspaceID &&
		this.UserID == other.UserID && this.StartDate == other.StartDate &&
		this.EndDate == other.EndDate && this.Cancelled == other.Cancelled &&
		this.CreatedBy == other.CreatedBy && this.WorkspaceName == other.WorkspaceID &&
		this.UserName == other.UserName && this.FloorID == other.FloorID && this.FloorName == other.FloorName
}

type Assignment struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_time"`
	EndDate     time.Time `json:"end_time"`
}
