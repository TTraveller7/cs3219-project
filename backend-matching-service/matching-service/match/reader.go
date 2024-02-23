package match

import (
	"context"
	"encoding/json"

	"github.com/backend-common/common"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	topic  string
}

var (
	consumers        [3]*Consumer
	consumersContext context.Context
	cancelConsumers  func()
	quit             chan int
)

func InitConsumers(q chan int) {
	consumers[0] = getConsumer(common.DIFFICULTY_EASY)
	consumers[1] = getConsumer(common.DIFFICULTY_MEDIUM)
	consumers[2] = getConsumer(common.DIFFICULTY_HARD)
	quit = q
}

func getConsumer(difficulty string) *Consumer {
	return &Consumer{
		reader: initReader(difficulty),
		topic:  difficulty,
	}
}

func initReader(topic string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Topic:   topic,
		Brokers: []string{getKafkaUrl()},
	})
	return r
}

func StartConsumers(mainContext context.Context) {
	consumersContext, cancelConsumers = context.WithCancel(mainContext)
	for i := range consumers {
		go consumers[i].start(consumersContext)
	}
}

func (c *Consumer) start(consumerContext context.Context) {
	log.Message("Consumer on topic " + c.topic + " starts")

	r := c.reader
	matchRequestCh := make(chan MatchRequest, 2)
	receiveMsgCtx, cancelReceiveMsg := context.WithCancel(consumerContext)
	defer cancelReceiveMsg()

	receiveMsg := func(ctx context.Context) {
		for {
			if ctx.Err() != nil {
				log.Message("ReceiveMsg routine canceled by parent context")
				return
			}

			msg, err := r.ReadMessage(ctx)
			if err != nil {
				log.Error("Fail to fetch message from consumer: ", err)
				continue
			}

			// deserialize user
			matchRequest := &MatchRequest{}
			if err := json.Unmarshal(msg.Value, matchRequest); err != nil {
				log.Error("Fail to deserialize messasge "+string(msg.Value), err)
				continue
			}

			isRequestValid, err := isMatchRequestValid(ctx, *matchRequest)
			if err != nil {
				log.Error("Fail to check if match request is valid:", err)
				continue
			}

			if isRequestValid {
				matchRequestCh <- *matchRequest
				return
			}
		}
	}

	var matchRequests [2]MatchRequest
	for {
		ptr := 0

		go receiveMsg(receiveMsgCtx)
		go receiveMsg(receiveMsgCtx)

		for areRequestsReady := false; !areRequestsReady; {
			select {
			case <-consumerContext.Done():
				r.Close()
				clearChannel(matchRequestCh)
				log.Message("Consumer on topic " + c.topic + " closes")
				return

			case matchRequest := <-matchRequestCh:
				matchRequests[ptr] = matchRequest
				log.Message("Consumer received match request from topic " + c.topic)
				ptr++

				// Assume length of matchrequest is 2
				if ptr == 2 {
					ok, err := isMatchRequestValid(receiveMsgCtx, matchRequests[0])
					if !ok || err != nil {
						matchRequests[0] = matchRequests[1]
						ptr--
						go receiveMsg(receiveMsgCtx)
					}
				}

				areRequestsReady = ptr >= len(matchRequests)
			}
		}

		clearChannel(matchRequestCh)

		difficulty := c.topic
		usernameA := matchRequests[0].Username
		usernameB := matchRequests[1].Username

		// Save match
		match := newMatch(difficulty, usernameA, usernameB)

		if err := saveMatch(consumerContext, match); err != nil {
			quit <- 0
			return
		} else {
			log.Message("Match created: " + match.toString())
		}

		// update user
		updateUserStatus(consumerContext, usernameA, match.MatchId, common.USER_STATUS_MATCHED)
		updateUserStatus(consumerContext, usernameB, match.MatchId, common.USER_STATUS_MATCHED)
		log.Message("User status updated")
	}
}

func CloseConsumers() {
	cancelConsumers()
}

func clearChannel[T any](ch chan T) {
	for len(ch) > 0 {
		<-ch
	}
}
