package import_values

import (
	"ValuesImporter/domain/entity"
	mqtt_transport "ValuesImporter/external/mqtt-transport"
	"ValuesImporter/persistence/mariadb/repository"
	save_value "ValuesImporter/usecase/save-value"
	"github.com/shettyh/threadpool"
	"github.com/tidwall/gjson"
	"strconv"
)

type ProcessorTask struct {
	topic     string
	message   []byte
	instances []entity.EntityParameterInstance
	importer  *EntityParameterValuesThreadedImporter
}

type EntityParameterValuesThreadedImporter struct {
	dataSourceRepository              repository.DataSourceRepository
	entityParameterInstanceRepository repository.EntityParameterInstanceRepository
	saveValue                         save_value.SaveValue
	mqttInstances                     []mqtt_transport.PahoMqtt
}

func NewEntityParameterValuesThreadedImporter(
	dataSourceRepository repository.DataSourceRepository,
	entityParameterInstanceRepository repository.EntityParameterInstanceRepository,
	saveValue save_value.SaveValue) *EntityParameterValuesThreadedImporter {
	return &EntityParameterValuesThreadedImporter{dataSourceRepository: dataSourceRepository, entityParameterInstanceRepository: entityParameterInstanceRepository, saveValue: saveValue}
}

func (importer *EntityParameterValuesThreadedImporter) Start() {
	var dataSources []entity.DataSource = importer.dataSourceRepository.GetDataSources()

	config := make(map[int]configurations, 0)

	pool := threadpool.NewThreadPool(200, 1000000)

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
			task := &ProcessorTask{topic: topic, message: message, instances: topicsMap[topic], importer: importer}
			pool.Execute(task)
		})

		importer.mqttInstances = append(importer.mqttInstances, *mqtt)
	}

}

func (t *ProcessorTask) Run() {
	for _, instance := range t.instances {
		if instance.LiveVariable != nil && *instance.LiveVariable != "" {
			var timestamp string = ""
			if instance.LiveTimestampVariable != nil && *instance.LiveTimestampVariable != "" {
				timestamp = gjson.Get(string(t.message), *instance.LiveTimestampVariable).String()
			}
			t.importer.saveValue.Save(instance, gjson.Get(string(t.message), *instance.LiveVariable).String(), timestamp)
		}
	}
}
