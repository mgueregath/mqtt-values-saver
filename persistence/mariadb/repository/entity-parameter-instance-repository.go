package repository

import (
	"ValuesImporter/domain/entity"
	"ValuesImporter/persistence"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type EntityParameterInstanceRepository struct {
	poolConnection *persistence.PoolConnection
}

var entityParameterInstanceRepository *EntityParameterInstanceRepository

func NewEntityParameterInstanceRepository(poolConnection *persistence.PoolConnection) *EntityParameterInstanceRepository {
	entityParameterInstanceRepository := EntityParameterInstanceRepository{poolConnection: poolConnection}
	return &entityParameterInstanceRepository
}

func (repository *EntityParameterInstanceRepository) GetEntityParameterInstances() []entity.EntityParameterInstance {
	var instances []entity.EntityParameterInstance
	err := repository.poolConnection.Query(
		fmt.Sprintf("SELECT %s FROM entity_parameter_instance", sqlstruct.Columns(entity.EntityParameterInstance{})),
		&instances)

	if err != nil {
		panic(err)
	}

	return instances
}
