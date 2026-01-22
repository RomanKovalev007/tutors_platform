package kafka

type KafkaConfig struct {
	Brokers string `env:"KAFKA_BROKERS" envDefault:"kafka:9092"`
	Topic   string `env:"KAFKA_TOPIC"`
	GroupID string `env:"KAFKA_GROUP_ID"`
}