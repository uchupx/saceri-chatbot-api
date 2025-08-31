package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type QueueMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, queueName string, message *QueueMessage) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return q.client.RPush(ctx, queueName, msgBytes).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (*QueueMessage, error) {
	result, err := q.client.BLPop(ctx, timeout, queueName).Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, redis.Nil
	}

	msgBytes := []byte(result[1])
	var message QueueMessage

	if err := json.Unmarshal(msgBytes, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
