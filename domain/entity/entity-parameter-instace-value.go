package entity

import "time"

type EntityParameterInstanceValue struct {
	Id                int64      `sql:"piv_id"`
	Timestamp         *time.Time `sql:"piv_timestamp"`
	TimestampString   *string    `sql:"piv_timestamp_string"`
	Saved             *time.Time `sql:"piv_saved"`
	Value             *float64   `sql:"piv_value"`
	ValueString       *string    `sql:"piv_timestamp_string"`
	ParameterInstance int        `sql:"piv_parameter_instance"`
}
