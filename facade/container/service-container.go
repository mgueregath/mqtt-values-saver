package container

import (
	"ValuesImporter/facade/environment"
	"ValuesImporter/persistence"
	"ValuesImporter/persistence/mariadb"
	"ValuesImporter/persistence/mariadb/repository"
	import_values "ValuesImporter/usecase/import-values"
	save_value "ValuesImporter/usecase/save-value"
)

var ENVIRONMENT *environment.Environment

func CreateApp(ENV *environment.Environment) {
	ENVIRONMENT = ENV
	startPersistence()
	createUsecases()
	startWorkers()

	for {

	}
}

func startPersistence() {
	poolConnection := mariadb.NewMariaDB(ENVIRONMENT)
	createRepositories(poolConnection)

}

var importerWorker *import_values.EntityParameterValuesImporter
var importerThreadedWorker *import_values.EntityParameterValuesThreadedImporter

func startWorkers() {
	//importerWorker = import_values.NewEntityParameterValuesImporter(*dataSourceRepository, *entityParameterInstanceRepository, *saveValue)
	importerThreadedWorker = import_values.NewEntityParameterValuesThreadedImporter(*dataSourceRepository, *entityParameterInstanceRepository, *saveValue)
	importerThreadedWorker.Start()

}

var dataSourceRepository *repository.DataSourceRepository
var entityParameterInstanceRepository *repository.EntityParameterInstanceRepository
var entityParameterInstanceValueRepository *repository.EntityParameterInstanceValueRepository

func createRepositories(poolConnection *persistence.PoolConnection) {
	dataSourceRepository = repository.NewDataSourceRepository(poolConnection)
	entityParameterInstanceRepository = repository.NewEntityParameterInstanceRepository(poolConnection)
	entityParameterInstanceValueRepository = repository.NewEntityParameterInstanceValueRepository(poolConnection)
}

var saveValue *save_value.SaveValue

func createUsecases() {
	saveValue = save_value.NewSaveValue(*entityParameterInstanceValueRepository)
}
