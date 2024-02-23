package match

import (
	"context"

	"github.com/backend-common/common"
	"github.com/segmentio/kafka-go"
)

var (
	easyWriter   *kafka.Writer
	mediumWriter *kafka.Writer
	hardWriter   *kafka.Writer

	writersCtx      context.Context
	cancelWritesCtx func()
)

func initWriters(ctx context.Context) {
	easyWriter = initWriter(common.DIFFICULTY_EASY)
	mediumWriter = initWriter(common.DIFFICULTY_MEDIUM)
	hardWriter = initWriter(common.DIFFICULTY_HARD)

	writersCtx, cancelWritesCtx = context.WithCancel(ctx)
}

func initWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(getKafkaUrl()),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		BatchSize:              2,
		BatchTimeout:           1,
		RequiredAcks:           1,
	}
	return w
}

func closeWriters() {
	closeWriter(easyWriter, mediumWriter, hardWriter)
	cancelWritesCtx()
}

func closeWriter(writers ...*kafka.Writer) {
	for i := range writers {
		if err := writers[i].Close(); err != nil {
			log.Error("Fail to close kafka writer:", err)
		}
	}
}

func WriteToEasyQueue(msg []byte) bool {
	return write(msg, easyWriter)
}

func WriteToMediumQueue(msg []byte) bool {
	return write(msg, mediumWriter)
}

func WriteToHardQueue(msg []byte) bool {
	return write(msg, hardWriter)
}

func write(msg []byte, w *kafka.Writer) bool {
	ctx, cancelFunc := context.WithCancel(writersCtx)
	defer cancelFunc()

	err := w.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		log.Error("Fail to write message "+string(msg)+" to "+w.Topic, err)
		return false
	}
	return true
}
