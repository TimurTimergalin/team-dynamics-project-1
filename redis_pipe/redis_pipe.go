package redis_pipe

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

package redislib

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type Config struct {
	Addr                   string
	Password               string
	DB                     int
	MaxRetries             int
	MinRetryBackoff        time.Duration
	MaxRetryBackoff        time.Duration
	CircuitBreakerName     string
	CircuitMaxRequests     uint32
	CircuitTimeout         time.Duration
	CircuitTripThreshold   uint32
	InChannelSize          int
	OutChannelSize         int
	ConnectTimeout         time.Duration
	MaxReconnectBackoff    time.Duration
	MaxReconnectAttempts   int
	PingInterval           time.Duration
	OperationRetryCount    int
}

func DefaultConfig() *Config {
	return &Config{
		Addr:                   "localhost:6379",
		DB:                     0,
		MaxRetries:             5,
		MinRetryBackoff:        8 * time.Millisecond,
		MaxRetryBackoff:        512 * time.Millisecond,
		CircuitBreakerName:     "redis",
		CircuitMaxRequests:     1,
		CircuitTimeout:         5 * time.Second,
		CircuitTripThreshold:   3,
		InChannelSize:          100,
		OutChannelSize:         100,
		ConnectTimeout:         5 * time.Second,
		MaxReconnectBackoff:    30 * time.Second,
		MaxReconnectAttempts:   10,
		PingInterval:           10 * time.Second,
		OperationRetryCount:    3,
	}
}

type EventSource struct {
	Channel     *string
	OperationId *string
}

type Event struct {
	Type    string
	Source  EventSource
	Payload []byte
	Err     error
}

type Operation interface {
	GetId() string
	Execute(r *Redis) error
}

type baseOperation struct {
	id string
}

func (b *baseOperation) GetId() string {
	return b.id
}

type subscribeOp struct {
	baseOperation
	channel string
}

func newSubscribeOp(channel string) *subscribeOp {
	return &subscribeOp{
		baseOperation: baseOperation{id: uuid.New().String()},
		channel:       channel,
	}
}

func (op *subscribeOp) Execute(r *Redis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subs[op.channel] {
		return nil
	}
	if r.pubsub != nil {
		if err := r.pubsub.Subscribe(context.Background(), op.channel); err != nil {
			return err
		}
	}
	r.subs[op.channel] = true
	return nil
}

type unsubscribeOp struct {
	baseOperation
	channel string
}

func newUnsubscribeOp(channel string) *unsubscribeOp {
	return &unsubscribeOp{
		baseOperation: baseOperation{id: uuid.New().String()},
		channel:       channel,
	}
}

func (op *unsubscribeOp) Execute(r *Redis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.subs[op.channel] {
		return nil
	}
	if r.pubsub != nil {
		if err := r.pubsub.Unsubscribe(context.Background(), op.channel); err != nil {
			return err
		}
	}
	delete(r.subs, op.channel)
	return nil
}

type customOp struct {
	baseOperation
	fn func(*redis.Client) error
}

func (op *customOp) Execute(r *Redis) error {
	r.mu.RLock()
	client := r.client
	r.mu.RUnlock()
	if client == nil {
		return errors.New("redis client not available")
	}
	return op.fn(client)
}

func MakeOp(fn func(*redis.Client) error) Operation {
	return &customOp{
		baseOperation: baseOperation{id: uuid.New().String()},
		fn:            fn,
	}
}

type Redis struct {
	client     *redis.Client
	pubsub     *redis.PubSub
	config     *Config
	breaker    *gobreaker.CircuitBreaker
	in         chan Operation
	out        chan Event
	subs       map[string]bool
	mu         sync.RWMutex
	done       chan struct{}
	wg         sync.WaitGroup
}

func NewRedis(config *Config) *Redis {
	if config == nil {
		config = DefaultConfig()
	}
	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        config.CircuitBreakerName,
		MaxRequests: config.CircuitMaxRequests,
		Timeout:     config.CircuitTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > config.CircuitTripThreshold
		},
	})
	return &Redis{
		config:  config,
		breaker: breaker,
		in:      make(chan Operation, config.InChannelSize),
		out:     make(chan Event, config.OutChannelSize),
		subs:    make(map[string]bool),
		mu:      sync.RWMutex{},
		done:    make(chan struct{}),
	}
}

func (r *Redis) Subscribe(channel string) {
	r.in <- newSubscribeOp(channel)
}

func (r *Redis) Unsubscribe(channel string) {
	r.in <- newUnsubscribeOp(channel)
}

func (r *Redis) Run() {
	r.wg.Add(1)
	defer r.wg.Done()
	r.runLoop()
}

func (r *Redis) runLoop() {
	psCh, err := r.connectWithRetry()
	if err != nil {
		r.out <- Event{Type: "unavailable_permanent", Source: EventSource{}, Err: err}
		return
	}

	pingTicker := time.NewTicker(r.config.PingInterval)
	defer pingTicker.Stop()

	for {
		select {
		case <-r.done:
			return

		case <-pingTicker.C:
			r.in <- MakeOp(func(client *redis.Client) error {
				ctx, cancel := context.WithTimeout(context.Background(), r.config.ConnectTimeout)
				defer cancel()
				return client.Ping(ctx).Err()
			})

		case op, ok := <-r.in:
			if !ok {
				return
			}

			opId := op.GetId()
			var err error

			for attempt := 1; attempt <= r.config.OperationRetryCount; attempt++ {
				err = op.Execute(r)
				if err == nil {
					break
				}

				if r.shouldReconnectOnError(err) {
					r.out <- Event{
						Type:   "error",
						Source: EventSource{OperationId: &opId},
						Err:    fmt.Errorf("operation failed (attempt %d/%d): %w", attempt, r.config.OperationRetryCount, err),
					}

					go r.triggerReconnect()
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				} else {
					r.out <- Event{
						Type:   "error",
						Source: EventSource{OperationId: &opId},
						Err:    fmt.Errorf("operation failed: %w", err),
					}
					break
				}
			}

		case msg, ok := <-psCh:
			if !ok {
				r.out <- Event{Type: "error", Source: EventSource{}, Err: errors.New("pubsub connection lost")}
				newPsCh, err := r.reconnectWithRetry()
				if err != nil {
					r.out <- Event{Type: "unavailable_permanent", Source: EventSource{}, Err: err}
					return
				}
				psCh = newPsCh
				continue
			}

			channel := msg.Channel
			r.out <- Event{
				Type:    "message",
				Source:  EventSource{Channel: &channel},
				Payload: []byte(msg.Payload),
			}
		}
	}
}

func (r *Redis) shouldReconnectOnError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "reset by peer") ||
		strings.Contains(msg, "EOF") ||
		errors.Is(err, redis.ErrClosed)
}

func (r *Redis) triggerReconnect() {
	select {
	case r.in <- MakeOp(func(client *redis.Client) error {
		return nil // No-op, just triggers reconnection check
	}):
	default:
	}
}

func (r *Redis) connectWithRetry() (<-chan *redis.Message, error) {
	attempt := 1
	maxAttempts := r.config.MaxReconnectAttempts

	for {
		if maxAttempts > 0 && attempt > maxAttempts {
			return nil, fmt.Errorf("failed to connect after %d attempts", maxAttempts)
		}

		select {
		case <-r.done:
			return nil, errors.New("shutdown during connection")
		default:
		}

		err := r.connect()
		if err == nil {
			return r.pubsub.Channel(), nil
		}

		r.out <- Event{Type: "unavailable_temporary", Source: EventSource{}, Err: err}

		backoff := time.Duration(1<<uint(attempt-1)) * time.Second
		if backoff > r.config.MaxReconnectBackoff {
			backoff = r.config.MaxReconnectBackoff
		}

		time.Sleep(backoff)
		attempt++
	}
}

func (r *Redis) reconnectWithRetry() (<-chan *redis.Message, error) {
	r.mu.Lock()
	if r.pubsub != nil {
		_ = r.pubsub.Close()
		r.pubsub = nil
	}
	r.mu.Unlock()
	return r.connectWithRetry()
}

func (r *Redis) connect() error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		r.mu.Lock()
		defer r.mu.Unlock()

		if r.client != nil {
			_ = r.client.Close()
			r.client = nil
		}

		client := redis.NewClient(&redis.Options{
			Addr:            r.config.Addr,
			Password:        r.config.Password,
			DB:              r.config.DB,
			MaxRetries:      r.config.MaxRetries,
			MinRetryBackoff: r.config.MinRetryBackoff,
			MaxRetryBackoff: r.config.MaxRetryBackoff,
		})

		ctx, cancel := context.WithTimeout(context.Background(), r.config.ConnectTimeout)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			return nil, err
		}

		r.client = client

		pubsub := client.Subscribe(ctx)
		channels := make([]string, 0, len(r.subs))
		for ch := range r.subs {
			channels = append(channels, ch)
		}

		if len(channels) > 0 {
			if err := pubsub.Subscribe(ctx, channels...); err != nil {
				return nil, err
			}
		}

		r.pubsub = pubsub
		return nil, nil
	})

	return err
}

func (r *Redis) Close() error {
	close(r.done)
	r.wg.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pubsub != nil {
		_ = r.pubsub.Close()
		r.pubsub = nil
	}

	if r.client != nil {
		err := r.client.Close()
		r.client = nil
		return err
	}

	return nil
}

func (r *Redis) Out() <-chan Event {
	return r.out
}

func (r *Redis) In() chan<- Operation {
	return r.in
}