package amqp

import (
	"fmt"
	"testing"
	"time"

	"github.com/printesoi/machinery/v2/brokers/iface"
	"github.com/printesoi/machinery/v2/config"
	"github.com/printesoi/machinery/v2/tasks"
	amqp "github.com/rabbitmq/amqp091-go"
)

type doNothingProcessor struct{}

func (_ doNothingProcessor) Process(signature *tasks.Signature) error {
	return fmt.Errorf("failed")
}

func (_ doNothingProcessor) CustomQueue() string {
	return "oops"
}

func (_ doNothingProcessor) PreConsumeHandler() bool {
	return true
}

func TestConsume(t *testing.T) {
	var (
		iBroker    iface.Broker
		deliveries = make(chan amqp.Delivery, 3)
		closeChan  chan *amqp.Error
		processor  doNothingProcessor
	)

	t.Run("with deliveries more than the number of concurrency", func(t *testing.T) {
		iBroker = New(&config.Config{})
		broker, _ := iBroker.(*Broker)
		errChan := make(chan error)

		// simulate that there are too much deliveries
		go func() {
			for i := 0; i < 3; i++ {
				deliveries <- amqp.Delivery{} // broker.consumeOne() will complain this error: Received an empty message
			}
		}()

		go func() {
			err := broker.consume(deliveries, 2, processor, closeChan)
			if err != nil {
				errChan <- err
			}
		}()

		select {
		case <-errChan:
		case <-time.After(1 * time.Second):
			t.Error("Maybe deadlock")
		}
	})
}
