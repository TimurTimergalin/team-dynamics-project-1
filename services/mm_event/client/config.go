package client

import "time"

type Config struct {
	MessageReceivedTimeout time.Duration
	CheckMatchPeriod       time.Duration
	CheckInPoolPeriod      time.Duration
	HubRegisterTimeout     time.Duration
	ConnectionTTL          time.Duration
}
