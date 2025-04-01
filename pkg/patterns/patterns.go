package main

import (
	"fmt"
	"log"
	"time"
)

type Operation func() error

func Retry(effector Operation, retries int, delay time.Duration) Operation {
	return func() error {
		for r := 0; r < retries; r++ {
			err := effector()
			if err == nil {
				return nil
			}
			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)
			time.Sleep(delay)
			delay *= 2
		}

		return effector()
	}
}

func Timeout(operation Operation, timeout time.Duration) Operation {
	return func() error {
		done := make(chan error, 1)
		go func() {
			done <- operation()
		}()
		select {
		case err := <-done:
			return err
		case <-time.After(timeout):
			return fmt.Errorf("timeout")
		}
	}
}

type Process func(msg string) error
type Messages []string

type DeadLetterQueue struct {
	messages Messages
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{messages: make(Messages, 0)}
}

func (dlq *DeadLetterQueue) GetMessages() Messages {
	return dlq.messages
}

func (dlq *DeadLetterQueue) AddMessage(message string) {
	dlq.messages = append(dlq.messages, message)
}

func ProcessWithDLQ(messages Messages, process Process, dlq *DeadLetterQueue) {
	for _, msg := range messages {
		err := process(msg)
		if err != nil {
			dlq.AddMessage(msg)
		}
	}
}
