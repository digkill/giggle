package event

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/nats.go"
	"log"

	"github.com/digkill/giggle/schema"
)

type NatsEventStore struct {
	nc                        *nats.Conn
	giggleCreatedSubscription *nats.Subscription
	giggleCreatedChan         chan GiggleCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (es *NatsEventStore) SubscribeGiggleCreated() (<-chan GiggleCreatedMessage, error) {
	m := GiggleCreatedMessage{}
	es.giggleCreatedChan = make(chan GiggleCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	es.giggleCreatedSubscription, err = es.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				if err := es.readMessage(msg.Data, &m); err != nil {
					log.Fatal(err)
				}
				es.giggleCreatedChan <- m
			}
		}
	}()
	return (<-chan GiggleCreatedMessage)(es.giggleCreatedChan), nil
}

func (es *NatsEventStore) OnGiggleCreated(f func(GiggleCreatedMessage)) (err error) {
	m := GiggleCreatedMessage{}
	es.giggleCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		if err := es.readMessage(msg.Data, &m); err != nil {
			log.Fatal(err)
		}
		f(m)
	})
	return
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.giggleCreatedSubscription != nil {
		if err := es.giggleCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	close(es.giggleCreatedChan)
}

func (es *NatsEventStore) PublishGiggleCreated(giggle schema.Giggle) error {
	m := GiggleCreatedMessage{giggle.ID, giggle.Body, giggle.CreatedAt}
	data, err := es.writeMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
