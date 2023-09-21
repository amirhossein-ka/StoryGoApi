// Code generated by ent, DO NOT EDIT.

package ent

import (
	"StoryGoAPI/ent/guestuser"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// GuestUser is the model entity for the GuestUser schema.
type GuestUser struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Token holds the value of the "token" field.
	Token string `json:"token,omitempty"`
	// VersionNumber holds the value of the "version_number" field.
	VersionNumber int `json:"version_number,omitempty"`
	// OperationSystem holds the value of the "operation_system" field.
	OperationSystem string `json:"operation_system,omitempty"`
	// UserAgent holds the value of the "user_agent" field.
	UserAgent string `json:"user_agent,omitempty"`
	// DisplayDetails holds the value of the "display_details" field.
	DisplayDetails string `json:"display_details,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the GuestUserQuery when eager-loading is set.
	Edges        GuestUserEdges `json:"edges"`
	selectValues sql.SelectValues
}

// GuestUserEdges holds the relations/edges for other nodes in the graph.
type GuestUserEdges struct {
	// Followed holds the value of the followed edge.
	Followed []*User `json:"followed,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// FollowedOrErr returns the Followed value or an error if the edge
// was not loaded in eager-loading.
func (e GuestUserEdges) FollowedOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Followed, nil
	}
	return nil, &NotLoadedError{edge: "followed"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*GuestUser) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case guestuser.FieldID, guestuser.FieldVersionNumber:
			values[i] = new(sql.NullInt64)
		case guestuser.FieldToken, guestuser.FieldOperationSystem, guestuser.FieldUserAgent, guestuser.FieldDisplayDetails:
			values[i] = new(sql.NullString)
		case guestuser.FieldCreatedAt, guestuser.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the GuestUser fields.
func (gu *GuestUser) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case guestuser.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			gu.ID = int(value.Int64)
		case guestuser.FieldToken:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field token", values[i])
			} else if value.Valid {
				gu.Token = value.String
			}
		case guestuser.FieldVersionNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field version_number", values[i])
			} else if value.Valid {
				gu.VersionNumber = int(value.Int64)
			}
		case guestuser.FieldOperationSystem:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field operation_system", values[i])
			} else if value.Valid {
				gu.OperationSystem = value.String
			}
		case guestuser.FieldUserAgent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_agent", values[i])
			} else if value.Valid {
				gu.UserAgent = value.String
			}
		case guestuser.FieldDisplayDetails:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display_details", values[i])
			} else if value.Valid {
				gu.DisplayDetails = value.String
			}
		case guestuser.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				gu.CreatedAt = value.Time
			}
		case guestuser.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				gu.DeletedAt = new(time.Time)
				*gu.DeletedAt = value.Time
			}
		default:
			gu.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the GuestUser.
// This includes values selected through modifiers, order, etc.
func (gu *GuestUser) Value(name string) (ent.Value, error) {
	return gu.selectValues.Get(name)
}

// QueryFollowed queries the "followed" edge of the GuestUser entity.
func (gu *GuestUser) QueryFollowed() *UserQuery {
	return NewGuestUserClient(gu.config).QueryFollowed(gu)
}

// Update returns a builder for updating this GuestUser.
// Note that you need to call GuestUser.Unwrap() before calling this method if this GuestUser
// was returned from a transaction, and the transaction was committed or rolled back.
func (gu *GuestUser) Update() *GuestUserUpdateOne {
	return NewGuestUserClient(gu.config).UpdateOne(gu)
}

// Unwrap unwraps the GuestUser entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (gu *GuestUser) Unwrap() *GuestUser {
	_tx, ok := gu.config.driver.(*txDriver)
	if !ok {
		panic("ent: GuestUser is not a transactional entity")
	}
	gu.config.driver = _tx.drv
	return gu
}

// String implements the fmt.Stringer.
func (gu *GuestUser) String() string {
	var builder strings.Builder
	builder.WriteString("GuestUser(")
	builder.WriteString(fmt.Sprintf("id=%v, ", gu.ID))
	builder.WriteString("token=")
	builder.WriteString(gu.Token)
	builder.WriteString(", ")
	builder.WriteString("version_number=")
	builder.WriteString(fmt.Sprintf("%v", gu.VersionNumber))
	builder.WriteString(", ")
	builder.WriteString("operation_system=")
	builder.WriteString(gu.OperationSystem)
	builder.WriteString(", ")
	builder.WriteString("user_agent=")
	builder.WriteString(gu.UserAgent)
	builder.WriteString(", ")
	builder.WriteString("display_details=")
	builder.WriteString(gu.DisplayDetails)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(gu.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := gu.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// GuestUsers is a parsable slice of GuestUser.
type GuestUsers []*GuestUser
