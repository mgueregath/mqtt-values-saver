package repository

import (
	"ValuesImporter/domain/entity"
	"ValuesImporter/persistence"
)

type EntityParameterInstanceValueRepository struct {
	poolConnection *persistence.PoolConnection
}

var entityParameterInstanceValueRepository *EntityParameterInstanceValueRepository

func NewEntityParameterInstanceValueRepository(poolConnection *persistence.PoolConnection) *EntityParameterInstanceValueRepository {
	if entityParameterInstanceValueRepository == nil {
		entityParameterInstanceValueRepository = &EntityParameterInstanceValueRepository{poolConnection: poolConnection}

	}
	return entityParameterInstanceValueRepository
}

func (repository *EntityParameterInstanceValueRepository) SaveEntityParameterValueInstances(value entity.EntityParameterInstanceValue) {
	err := repository.poolConnection.Save(
		"INSERT INTO `entity_parameter_instance_value` (`piv_timestamp`, `piv_value`, `piv_timestamp_string`, `piv_value_string`, `piv_parameter_instance`, `piv_saved`) VALUES('%s', %f, '%s', '%s', %d, NOW())", []any{*value.TimestampString, *value.Value, *value.TimestampString, *value.ValueString, value.ParameterInstance})

	if err != nil {
		panic(err)
	}
}
