package save_value

import (
	"ValuesImporter/domain/entity"
	"ValuesImporter/persistence/mariadb/repository"
	"strconv"
	"time"
)

type SaveValue struct {
	entityParameterInstanceValueRepository repository.EntityParameterInstanceValueRepository
}

type ISaveValue interface {
	Save(entityParameterInstance entity.EntityParameterInstance, value string, timestamp string) error
}

func NewSaveValue(
	entityParameterInstanceValueRepository repository.EntityParameterInstanceValueRepository) *SaveValue {
	return &SaveValue{entityParameterInstanceValueRepository}
}

func (s *SaveValue) Save(entityParameterInstance entity.EntityParameterInstance, value string, timestamp string) error {
	valueNumber, _ := strconv.ParseFloat(value, 64)
	var timestampDate time.Time = time.Now()
	if timestamp != "" && &timestamp != nil {
		timestampDate, _ = time.Parse("2006-01-02 15:04:05", timestamp)
	} else {
		timestamp = timestampDate.Format("2006-01-02 15:04:05")
	}

	valueObj := entity.EntityParameterInstanceValue{
		ParameterInstance: entityParameterInstance.Id,
		Value:             &valueNumber,
		ValueString:       &value,
		Timestamp:         &timestampDate,
		TimestampString:   &timestamp,
	}
	s.entityParameterInstanceValueRepository.SaveEntityParameterValueInstances(valueObj)

	return nil
}
