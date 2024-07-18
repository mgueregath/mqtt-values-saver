package repository

import (
	"ValuesImporter/domain/entity"
	"ValuesImporter/persistence"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type DataSourceRepository struct {
	poolConnection *persistence.PoolConnection
}

var dataSourceRepository *DataSourceRepository

func NewDataSourceRepository(poolConnection *persistence.PoolConnection) *DataSourceRepository {
	dataSourceRepository := DataSourceRepository{poolConnection: poolConnection}
	return &dataSourceRepository
}

func (repository *DataSourceRepository) GetDataSources() []entity.DataSource {
	var dataSources []entity.DataSource
	err := repository.poolConnection.Query(
		fmt.Sprintf("SELECT %s FROM data_source", sqlstruct.Columns(entity.DataSource{})),
		&dataSources)

	if err != nil {
		panic(err)
	}

	return dataSources
}

func (repository *DataSourceRepository) GetMqttDataSources() []entity.DataSource {
	var dataSources []entity.DataSource
	err := repository.poolConnection.Query(
		fmt.Sprintf("SELECT %s FROM data_source WHERE dso_type = 1", sqlstruct.Columns(entity.DataSource{})),
		&dataSources)

	if err != nil {
		panic(err)
	}

	return dataSources
}
