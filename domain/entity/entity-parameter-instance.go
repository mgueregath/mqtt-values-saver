package entity

type EntityParameterInstance struct {
	Id                    int     `sql:"epi_id"`
	Hidden                *string `sql:"epi_hidden"`
	LiveVariable          *string `sql:"epi_live_variable"`
	Parameter             *int    `sql:"epi_parameter"`
	LiveDataSource        *int    `sql:"epi_live_data_source"`
	LiveTopic             *string `sql:"epi_live_topic"`
	LiveTimestampVariable *string `sql:"epi_live_timestamp_variable"`
}
