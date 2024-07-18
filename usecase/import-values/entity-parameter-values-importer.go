package import_values

import (
	"ValuesImporter/domain/entity"
	mqtt_transport "ValuesImporter/external/mqtt-transport"
	"ValuesImporter/persistence/mariadb/repository"
	save_value "ValuesImporter/usecase/save-value"
	"github.com/tidwall/gjson"
	"strconv"
)

type EntityParameterValuesImporter struct {
	dataSourceRepository              repository.DataSourceRepository
	entityParameterInstanceRepository repository.EntityParameterInstanceRepository
	saveValue                         save_value.SaveValue
	mqttInstances                     []mqtt_transport.PahoMqtt
}

func NewEntityParameterValuesImporter(
	dataSourceRepository repository.DataSourceRepository,
	entityParameterInstanceRepository repository.EntityParameterInstanceRepository,
	saveValue save_value.SaveValue) *EntityParameterValuesImporter {
	return &EntityParameterValuesImporter{dataSourceRepository: dataSourceRepository, entityParameterInstanceRepository: entityParameterInstanceRepository, saveValue: saveValue}
}

type configurations struct {
	dataSource entity.DataSource
	topics     []entity.EntityParameterInstance
}

func (importer *EntityParameterValuesImporter) Start() {
	var dataSources []entity.DataSource = importer.dataSourceRepository.GetDataSources()

	config := make(map[int]configurations, 0)

	var entityParameterInstances []entity.EntityParameterInstance = importer.entityParameterInstanceRepository.GetEntityParameterInstances()

	for _, ds := range dataSources {
		if ds.Type == 1 {
			config[ds.Id] = configurations{dataSource: ds, topics: make([]entity.EntityParameterInstance, 0)}
		}
	}

	for _, parameterInstance := range entityParameterInstances {
		if parameterInstance.LiveDataSource != nil && parameterInstance.LiveTopic != nil {
			id := *parameterInstance.LiveDataSource
			if ds, exist := config[id]; exist {
				ds.topics = append(ds.topics, parameterInstance)
				config[id] = ds
			}
		}
	}

	for _, config := range config {
		var mqtt = mqtt_transport.NewPahoMqtt(strconv.Itoa(config.dataSource.Id), *config.dataSource.Url, config.dataSource.Port, config.dataSource.Username, config.dataSource.Password)

		topics := make([]string, 0)

		topicsMap := make(map[string][]entity.EntityParameterInstance)

		for _, topic := range config.topics {
			topics = append(topics, *topic.LiveTopic)
			if _, exist := topicsMap[*topic.LiveTopic]; exist {
				topicsMap[*topic.LiveTopic] = []entity.EntityParameterInstance{}
			}
			topicsMap[*topic.LiveTopic] = append(topicsMap[*topic.LiveTopic], topic)
		}

		mqtt.Connect(topics, func(topic string, message []byte) {
			for _, instance := range topicsMap[topic] {
				importer.saveValue.Save(instance, gjson.Get(string(message), *instance.LiveVariable).String(), gjson.Get(string(message), *instance.LiveTimestampVariable).String())
			}
		})

		importer.mqttInstances = append(importer.mqttInstances, *mqtt)
	}

}
