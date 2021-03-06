package db

import (
	"errors"
	"go-api/model"
	"time"
)

var NotFoundError = errors.New("not found")
var EmptyError = errors.New("empty")

type DataStore struct {
	Closable
	WorkspaceProvider workspaceProvider
	BookingProvider   bookingProvider
	UserProvider      userProvider
	FloorProvider     floorProvider
	OfferingProvider  offeringProvider
	AssigneeProvider  assigneeProvider
}

type Closable interface {
	Close()
}

type workspaceProvider interface {
	GetOneWorkspace(id string) (*model.Workspace, error)
	UpdateWorkspace(id string, workspace *model.Workspace) error
	UpdateWorkspaceMetadata(id string, properties *model.Attrs) error
	CreateWorkspace(workspace *model.Workspace) (string, error)
	UpsertWorkspace(workspace *model.Workspace) (string, error)
	RemoveWorkspace(id string) error
	GetAllWorkspaces() ([]*model.Workspace, error)
	GetAllWorkspacesByFloor(floorId string) ([]*model.Workspace, error)
	FindAvailability(floorId string, start time.Time, end time.Time) ([]string, error)
	CountWorkspacesByFloor(floorId string) (int, error)
	CreateAssignment(userId, workspaceId string) error
	CreateAssignWorkspace(workspace *model.Workspace, userId string) (string, error)
	GetDeletedWorkspaces() ([]*model.Workspace, error)
	DeleteWorkspaces(ids []string) error
}

type bookingProvider interface {
	GetOneBooking(id string) (*model.Booking, error)
	GetOneExpandedBooking(id string) (*model.ExpandedBooking, error)
	GetAllBookings() ([]*model.Booking, error)
	GetAllExpandedBookings() ([]*model.ExpandedBooking, error)
	GetBookingsByWorkspaceID(id string) ([]*model.Booking, error)
	GetExpandedBookingsByWorkspaceID(id string) ([]*model.ExpandedBooking, error)
	GetBookingsByUserID(id string) ([]*model.Booking, error)
	GetExpandedBookingsByUserID(id string) ([]*model.ExpandedBooking, error)
	GetBookingsByDateRange(start time.Time, end time.Time) ([]*model.Booking, error)
	GetExpandedBookingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedBooking, error)
	CreateBooking(booking *model.Booking) (string, error)
	UpdateBooking(id string, booking *model.Booking) error
	RemoveBooking(id string) error
	GetExpiredBookings(since time.Time) ([]*model.Booking, error)
	DeleteBookings(ids []string) error
}

type userProvider interface {
	GetOneUser(id string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(user *model.User) error
	GetAssignedUsers(start, end time.Time) ([]*model.UserAssignment, error)
	GetAssignedUsersByTime(timestamp time.Time) ([]*model.UserAssignment, error)
	//UpdateUser(id string, user *model.User) error
	//RemoveUser(id string) error
}

type floorProvider interface {
	GetOneFloor(id string) (*model.Floor, error)
	GetAllFloors() ([]*model.Floor, error)
	GetAllFloorIDs() ([]string, error)
	CreateFloor(floor *model.Floor) (string, error)
	RemoveFloor(id string, force bool) error
	GetDeletedFloors() ([]*model.Floor, error)
	//UpdateFloor(id string, user *model.Floor) error
	DeleteFloors(ids []string) error
}

type offeringProvider interface {
	GetOneOffering(id string) (*model.Offering, error)
	GetOneExpandedOffering(id string) (*model.ExpandedOffering, error)
	GetAllOfferings() ([]*model.Offering, error)
	GetAllExpandedOfferings() ([]*model.ExpandedOffering, error)
	GetOfferingsByWorkspaceID(id string) ([]*model.Offering, error)
	GetExpandedOfferingsByWorkspaceID(id string) ([]*model.ExpandedOffering, error)
	GetOfferingsByUserID(id string) ([]*model.Offering, error)
	GetExpandedOfferingsByUserID(id string) ([]*model.ExpandedOffering, error)
	GetOfferingsByDateRange(start time.Time, end time.Time) ([]*model.Offering, error)
	GetExpandedOfferingsByDateRange(start time.Time, end time.Time) ([]*model.ExpandedOffering, error)
	GetOfferingsByWorkspaceIDAndDateRange(id string, start time.Time, end time.Time) (*model.Offering, error)
	CreateOffering(booking *model.Offering) (string, error)
	CreateDefaultOffering(booking *model.Offering) (string, error)
	UpdateOffering(id string, booking *model.Offering) error
	RemoveOffering(id string) error
	GetExpiredOfferings(since time.Time) ([]*model.Offering, error)
	DeleteOfferings(ids []string) error
}

type assigneeProvider interface {
	IsAssigned(id string, start time.Time, end time.Time) (bool, error)
	IsFullyAssigned(id string, start time.Time, end time.Time) (bool, error)
	GetExpiredAssignments(since time.Time) ([]*model.Assignment, error)
	DeleteAssignments(ids []string) error
}
