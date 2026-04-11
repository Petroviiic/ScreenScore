package main

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func GetRealClock() *RealClock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now().UTC()
}

type MockClock struct {
	FixedTime time.Time
}

func GetMockClock() *MockClock {
	return &MockClock{}
}

func (c *MockClock) Now() time.Time {
	return c.FixedTime
}
