package kite

import "koding/messaging/moh"

// HTTPMessenger is a struct that complies with the Messenger interface.
type HTTPMessenger struct {
	Subscriber *moh.Subscriber
	Requester  *moh.Requester

	// Messages published by the publisher but not consumber by Consume() yet.
	messages chan *[]byte
}

// NewHTTPMessenger returns a pointer to a new HTTPMessenger.
// Created HTTPMessenger will keep an open connection to the other side for
// consuming messages asynchronously.
func NewHTTPMessenger(kiteID string) (*HTTPMessenger, error) {
	const replierURL = "http://127.0.0.1:5556"
	const publisherURL = "ws://127.0.0.1:5557"

	requester, err := moh.NewRequester(replierURL)
	if err != nil {
		return nil, err
	}

	messages := make(chan *[]byte)
	handler := makeHandler(messages)
	subscriber, err := moh.NewSubscriber(publisherURL, handler)
	if err != nil {
		return nil, err
	}

	return &HTTPMessenger{
		Subscriber: subscriber,
		Requester:  requester,
		messages:   messages,
	}, nil
}

// Subscribe registers the filter for consuming messages.
func (h *HTTPMessenger) Subscribe(filter string) error {
	h.Subscriber.Subscribe(filter)
	return nil
}

// Unsubscribe unregisters the filter that is registered by Subscribe().
func (h *HTTPMessenger) Unsubscribe(filter string) error {
	h.Subscriber.Unsubscribe(filter)
	return nil
}

// Send is used for synchronous request/reply messaging pattern.
func (h *HTTPMessenger) Send(msg []byte) []byte {
	reply, err := h.Requester.Request(msg)
	if err != nil {
	}
	return reply
}

// Consume is a blocking function that is used for processing messages
// coming from Publisher.
func (h *HTTPMessenger) Consume(handler func([]byte)) {
	for msg := range h.messages {
		go handler(*msg)
	}
}

// makeHandler returns a function that queues messages to a channel
func makeHandler(messages chan<- *[]byte) func([]byte) {
	return func(msg []byte) {
		messages <- &msg
	}
}