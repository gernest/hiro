package bus

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/nsqio/go-nsq"
)

const (
	ScanTopic       = "scan"
	AccessTopic     = "access"
	PrintTopic      = "print"
	PrintStatsTopic = "print_progress"
)

type SilentLogger struct{}

func (SilentLogger) Output(calldepth int, s string) error { return nil }

// ScanEvent is an event fired when a qrcode embedded url is scanned.
type ScanEvent struct {
	// Status is the http status code of the scan request. eg 404, 200 etc.
	Status int `json:"status"`

	// ID is the uuid for the qrcode that is being scanned. This is taken directly
	// from the url /scan/:uuid
	ID string `json:"id"`

	// User is the uuid of the user who scanned the qrcode if any
	User string `json:"user"`

	Timestamp string `json:"timestamp"`
}

type AccessPolicy struct {
	Allowed []string `json:"allowed,omitempty"`
	Denied  []string `json:"denied,omitempty"`
}

// Producer this is a wrapper around nsq to allow sending  messages to the
// nsq service.
//
// All topics are published under the backend channel.
type Producer struct {
	producer *nsq.Producer
	Logger   *zap.Logger
	cancel   func()
}

// NewProducer returns a new *Producer instance that is connected to a producer
// nsqd.
//
//pings to  the daemon to ensure connections can be established.
func NewProducer(addr string) (*Producer, error) {
	cfg := nsq.NewConfig()
	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(); err != nil {
		return nil, err
	}
	p.SetLogger(SilentLogger{}, nsq.LogLevelInfo)
	b := &Producer{producer: p}
	return b, nil
}

// PublishScan sends the ScanEvent on the ScanTopic
func (b *Producer) PublishScan(s ScanEvent) {
	b.publish(ScanTopic, s)
}

// PublishAccessPolicy sends AccessPolicy message on AccessTopic.
func (b *Producer) PublishAccessPolicy(a AccessPolicy) {
	b.publish(AccessTopic, a)
}

func (b *Producer) publish(topic string, v interface{}) error {
	if b.Logger != nil {
		b.Logger.Info("publishing",
			zap.String("topic", topic),
		)
	}
	body, _ := json.Marshal(v)
	return b.producer.Publish(topic, body)
}
