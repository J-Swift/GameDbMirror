package model

import (
	"database/sql"
	"encoding/json"
)

type NullInt sql.NullInt32
type NullBool sql.NullBool
type NullString sql.NullString

func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}
func (ni *NullInt) UnmarshalJSON(data []byte) error {
	var obj int32
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*ni = NullInt{obj, true}
	return nil
}

func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var obj bool
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*nb = NullBool{obj, true}
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid || ns.String == "" {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var obj string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*ns = NullString{obj, true}
	return nil
}
