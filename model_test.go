package rapi

import "testing"

type M struct {
	Model
}

var m = M{}

func TestModelValidatePresence(t *testing.T) {
	m.ResetErrors()

	m.ValidatePresence("Name", "name")
	assertEqual(t, m.IsValid(), true)

	m.ValidatePresence("Name", "")
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateLength(t *testing.T) {
	m.ResetErrors()

	m.ValidateLength("Name", "name", -1, -1) // no limits
	assertEqual(t, m.IsValid(), true)

	m.ValidateLength("Name", "name", 4, -1) // min 4
	assertEqual(t, m.IsValid(), true)

	m.ValidateLength("Name", "name", -1, 4) // max 4
	assertEqual(t, m.IsValid(), true)

	m.ValidateLength("Name", "name", 4, 4) // min 4, max 4
	assertEqual(t, m.IsValid(), true)

	m.ValidateLength("Name", "name", 5, -1) // min 5
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateLength("Name", "name", -1, 3) // max 3
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateLength("Name", "name", 5, 30) // min 5, max 30
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateInt(t *testing.T) {
	m.ResetErrors()

	m.ValidateInt("Name", 10, -1, -1) // no limits
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt("Name", 10, 1, -1) // min 1
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt("Name", 10, -1, 10) // max 10
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt("Name", 10, 11, -1) // min 11
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateInt("Name", 10, -1, 9) // max 9
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateInt64(t *testing.T) {
	m.ResetErrors()

	m.ValidateInt64("Name", 10, -1, -1) // no limits
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt64("Name", 10, 1, -1) // min 1
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt64("Name", 10, -1, 10) // max 10
	assertEqual(t, m.IsValid(), true)

	m.ValidateInt64("Name", 10, 11, -1) // min 11
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateInt64("Name", 10, -1, 9) // max 9
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateFloat32(t *testing.T) {
	m.ResetErrors()

	m.ValidateFloat32("Name", 10.1, -1, -1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat32("Name", 10.1, 1, -1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat32("Name", 10.1, -1, 10.1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat32("Name", 10.1, 11, -1)
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateFloat32("Name", 10.1, -1, 9)
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateFloat64(t *testing.T) {
	m.ResetErrors()

	m.ValidateFloat64("Name", 10.1, -1, -1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat64("Name", 10.1, 1, -1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat64("Name", 10.1, -1, 10.1)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFloat64("Name", 10.1, 11, -1)
	assertEqual(t, m.IsValid(), false)

	m.ResetErrors()

	m.ValidateFloat64("Name", 10.1, -1, 9)
	assertEqual(t, m.IsValid(), false)
}

func TestModelValidateFormat(t *testing.T) {
	m.ResetErrors()

	m.ValidateFormat("IP", "1.1.1.1", `\A(\d{1,3}\.){3}\d{1,3}\z`)
	assertEqual(t, m.IsValid(), true)

	m.ValidateFormat("IP", "1.1.1", `\A(\d{1,3}\.){3}\d{1,3}\z`)
	assertEqual(t, m.IsValid(), false)
}
