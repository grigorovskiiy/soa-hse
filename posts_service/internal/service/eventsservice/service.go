package eventsservice

import (
	"context"
	"encoding/json"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/kafka"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	kafkaGo "github.com/segmentio/kafka-go"
)

type KafkaService struct {
	producer *kafka.BaseProducer
}

func NewKafkaService(cfg *config.Config, producer *kafka.BaseProducer) (*KafkaService, error) {
	client := kafkaGo.Client{
		Addr: kafkaGo.TCP(cfg.Brokers...),
	}

	topicsReq := CreateTopicsReq(cfg)
	_, err := client.CreateTopics(context.Background(), &topicsReq)
	if err != nil {
		logger.Logger.Error("error creating topics", "error", err.Error())
		return nil, err
	}

	return &KafkaService{producer: producer}, nil
}

func (s *KafkaService) SendEvent(ctx context.Context, topic string, upd any) error {
	msg, err := json.Marshal(upd)
	if err != nil {
		logger.Logger.Error("upd json marshal error", "error", err.Error())
		return err
	}

	if err := s.producer.Produce(ctx, topic, &kafkaGo.Message{
		Value: msg,
	}); err != nil {
		logger.Logger.Error("kafka produce error", "error", err.Error())
		return err
	}

	return nil
}

func CreateTopicsReq(cfg *config.Config) kafkaGo.CreateTopicsRequest {
	return kafkaGo.CreateTopicsRequest{
		Topics: []kafkaGo.TopicConfig{
			{
				Topic:             cfg.ViewsTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             cfg.CommentsTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             cfg.LikesTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	}
}
