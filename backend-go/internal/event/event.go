// Package event provides §1.6 EventBus in-memory implementation (doc v4.0).
package event

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Event is a single event (§1.6.2).
type Event struct {
	ID        string
	Topic     string
	Timestamp time.Time
	UserID    string
	Payload   interface{}
	Metadata  map[string]string
}

// EventHandler is called for each event (§1.6.2).
type EventHandler func(ctx context.Context, ev Event) error

// subEntry holds handler and id for removal.
type subEntry struct {
	id      int
	handler EventHandler
}

// Subscription represents a subscription (for Unsubscribe).
type Subscription struct {
	topic string
	id    int
}

// EventBus provides Publish and Subscribe (§1.6.2).
type EventBus struct {
	mu        sync.RWMutex
	subs      map[string][]subEntry
	idCounter int
}

// NewEventBus returns a new in-memory event bus.
func NewEventBus() *EventBus {
	return &EventBus{subs: make(map[string][]subEntry)}
}

// Publish sends an event to all subscribers of the topic.
func (b *EventBus) Publish(topic string, payload interface{}, userID string, metadata map[string]string) error {
	b.mu.RLock()
	entries := b.subs[topic]
	if len(entries) == 0 {
		b.mu.RUnlock()
		return nil
	}
	entries = append([]subEntry(nil), entries...)
	b.mu.RUnlock()

	b.idCounter++
	ev := Event{
		ID:        fmt.Sprintf("ev-%d", b.idCounter),
		Topic:     topic,
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Payload:   payload,
		Metadata:  metadata,
	}
	ctx := context.Background()
	for _, e := range entries {
		_ = e.handler(ctx, ev)
	}
	return nil
}

// Subscribe adds a handler for topic. Returns a Subscription for Unsubscribe.
func (b *EventBus) Subscribe(topic string, handler EventHandler) *Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.idCounter++
	id := b.idCounter
	b.subs[topic] = append(b.subs[topic], subEntry{id: id, handler: handler})
	return &Subscription{topic: topic, id: id}
}

// Unsubscribe removes the subscription.
func (b *EventBus) Unsubscribe(sub *Subscription) {
	if sub == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	list := b.subs[sub.topic]
	for i, e := range list {
		if e.id == sub.id {
			b.subs[sub.topic] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// Close is a no-op for in-memory bus (for interface compatibility).
func (b *EventBus) Close() error {
	return nil
}
