package consumer

import (
	"github.com/Shopify/sarama"
)

type CustomConsumer struct {
	ClaimFunc func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

func NewCustomConsumer(claimFunc func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error) *CustomConsumer {
	return &CustomConsumer{
		ClaimFunc: claimFunc,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *CustomConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *CustomConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *CustomConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return c.ClaimFunc(session, claim)
}
