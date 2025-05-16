package kafka

import (
	"context"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"hash/adler32"
)

type producer interface {
	WriteMessages(context.Context, ...kafka.Message) error
}

type BaseProducer struct {
	producer
}

func NewBaseProducer(lc fx.Lifecycle, cfg *config.Config) *BaseProducer {
	writer := &kafka.Writer{
		Addr:      kafka.TCP(cfg.Brokers...),
		Balancer:  &kafka.Hash{Hasher: adler32.New()},
		Transport: kafka.DefaultTransport,
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return writer.Close()
		},
	})

	return &BaseProducer{writer}
}

func (p *BaseProducer) Produce(ctx context.Context, topic string, msg *kafka.Message) error {
	msg.Topic = topic
	err := p.WriteMessages(ctx, *msg)

	if err != nil {
		logger.Logger.Error("error producing kafka message", "topic", topic, "error", err.Error())
	} else {
		logger.Logger.Info("kafka message produced successfully", "topic", topic)
	}

	return err
}
