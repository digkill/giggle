package event

import "github.com/digkill/giggle/schema"

type EventStore interface {
	Close()
	PublishGiggleCreated(giggle schema.Giggle) error
	SubscribeGiggleCreated() (<-chan GiggleCreatedMessage, error)
	OnGiggleCreated(f func(GiggleCreatedMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishGiggleCreated(giggle schema.Giggle) error {
	return impl.PublishGiggleCreated(giggle)
}

func SubscribeGiggleCreated() (<-chan GiggleCreatedMessage, error) {
	return impl.SubscribeGiggleCreated()
}

func OnGiggleCreated(f func(GiggleCreatedMessage)) error {
	return impl.OnGiggleCreated(f)
}
