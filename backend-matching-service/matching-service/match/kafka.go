package match

import (
	"context"
	"os"
	"strconv"

	"github.com/backend-common/common"
	"github.com/segmentio/kafka-go"
)

func SetupKafka() {
	kafkaSetupCtx, cancelKafkaSetup := context.WithCancel(mainContext)
	defer cancelKafkaSetup()
	for i := 0; i < KAFKA_RETRY_TIMES; i++ {
		select {
		case <-kafkaSetupCtx.Done():
			log.Message("Kafka setup is cancelled by parent context")
			return
		default:
			err := createTopics(kafkaSetupCtx)
			if err != nil {
				continue
			}

			initWriters(mainContext)
			return
		}
	}
	log.Message("Fail to connect to Kafka after " + strconv.Itoa(KAFKA_RETRY_TIMES) + " retries. Aboriting...")
	os.Exit(1)
}

func createTopics(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", getKafkaUrl())
	if err != nil {
		log.Error("Fail to connect to kafka:", err)
		return err
	}
	defer conn.Close()

	topicConfigs := [3]kafka.TopicConfig{}
	i := 0
	for difficulty := range common.DifficultySet() {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             difficulty,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	err = conn.CreateTopics(topicConfigs[0], topicConfigs[1], topicConfigs[2])
	if err != nil {
		log.Error("Fail to create topics:", err)
		return err
	}

	log.Message("Kafka topics created.")
	return nil
}

func CloseKafka() {
	closeWriters()
}

func TryConnectToKafka() {
	url := getKafkaUrl()
	topic := common.DIFFICULTY_EASY
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", url, topic, partition)
	defer conn.Close()

	if err != nil {
		log.Error("failed to dial leader:", err)
		os.Exit(1)
	}
}
