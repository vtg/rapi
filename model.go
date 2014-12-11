package rapi

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"
)

// ModelErrors errors type
type ModelErrors map[string][]string

// Model structure for base model
//	type User struct {
// 	    rapi.Model
// 	    Name string
// 	}
type Model struct {
	Id        int64       `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Errors    ModelErrors `sql:"-" json:"-"`
}

// ID returns ID of record
func (m *Model) ID() int64 {
	return m.Id
}

func (m *Model) ResetErrors() {
	m.Errors = make(ModelErrors)
}

// AddError adding error to record
func (m *Model) AddError(f string, t string) {
	if m.IsValid() {
		m.Errors = make(ModelErrors)
	}
	m.Errors[f] = append(m.Errors[f], t)
}

// IsValid returns true if no errors on record
func (m *Model) IsValid() bool {
	return len(m.Errors) == 0
}

// Valid placeholder for validation function
//  func (u *User) Valid() bool {
//  	u.ValidatePresence("Name", u.Name)
//  	return u.IsValid()
//  }
func (m *Model) Valid() bool {
	return m.IsValid()
}

// GetErrors returns record errors
func (m *Model) GetErrors() ModelErrors {
	return m.Errors
}

// SetErrors set record errors
func (m *Model) SetErrors(e ModelErrors) {
	m.Errors = e
}

// ValidatePresence validates string for presence
// 	m.ValidatePresence("Name", m.Name)
func (m *Model) ValidatePresence(f, v string) {
	if utf8.RuneCountInString(v) == 0 {
		m.AddError(f, "can't be blank")
	}
}

// ValidateLength validates string min, max length. -1 for any
// 	m.ValidateLength("password", m.Password, 6, 18) // min 6, max 18
func (m *Model) ValidateLength(f, v string, min, max int) {
	if min > 0 {
		if utf8.RuneCountInString(v) < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if utf8.RuneCountInString(v) > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateInt validates int min, max. -1 for any
// 	m.ValidateInt("number", 10, -1, 11)  // max 18
func (m *Model) ValidateInt(f string, v, min, max int) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateInt64 validates int64 min, max. -1 for any
// 	m.ValidateInt64("number", 10, 6, -1) // min 6
func (m *Model) ValidateInt64(f string, v, min, max int64) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFloat32 validates float32 min, max. -1 for any
// 	m.ValidateFloat32("number", 10.2, -1, 11)
func (m *Model) ValidateFloat32(f string, v, min, max float32) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFloat64 validates float64 min, max. -1 for any
// 	m.ValidateFloat64("number", 10.2, -1, 11)
func (m *Model) ValidateFloat64(f string, v, min, max float64) {
	if min > 0 {
		if v < min {
			m.AddError(f, fmt.Sprint("minimum length is", min))
		}
	}
	if max > 0 {
		if v > max {
			m.AddError(f, fmt.Sprint("maximum length is", max))
		}
	}
}

// ValidateFormat validates string format with regex string
// 	m.ValidateFormat("ip address", u.IP, `\A(\d{1,3}\.){3}\d{1,3}\z`)
func (m *Model) ValidateFormat(f, v, reg string) {
	if r, _ := regexp.MatchString(reg, v); !r {
		m.AddError(f, "invalid format")
	}
}

// BaseModel interface
type BaseModel interface {
	ID() int64
	Valid() bool

	AddError(string, string)
	SetErrors(ModelErrors)
	GetErrors() ModelErrors
	ResetErrors()
}
