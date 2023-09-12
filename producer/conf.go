package producer

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func getKafkaConf() *kafka.ConfigMap {
	topicUrl := os.Getenv("KAFKA_TOPIC_URL")
	username := os.Getenv("KAFKA_API_KEY")
	secret := os.Getenv("KAFKA_SECRET")

	return &kafka.ConfigMap{
		"bootstrap.servers": topicUrl,
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     username,
		"sasl.password":     secret,
	}
}
