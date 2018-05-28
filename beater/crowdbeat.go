package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/adamstauffer/crowdbeat/config"
)

type Crowdbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Crowdbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Crowdbeat) Run(b *beat.Beat) error {
	logp.Info("crowdbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		// process new audit events
		startDateTime := time.Now().Add(
			-bt.config.Period).Format("2006-01-02T15:04:05.999-0700")
		endDateTime := time.Now().Format("2006-01-02T15:04:05.999-0700")
		logp.Info("Requesting audit events between %v and %v",
			startDateTime, endDateTime)

		events := QueryAuditLog(bt, startDateTime, endDateTime)
		logp.Info("Received %v audit events", len(events))
		for _, event := range events {
			parsedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700",
				event.Timestamp)
			beatEvent := beat.Event{
				Timestamp: parsedTime,
				Fields: common.MapStr{
					"type":          b.Info.Name,
					"id":            event.Id,
					"author":        event.Author,
					"event_type":    event.EventType,
					"entities":      event.Entities,
					"ip_address":    event.IpAddress,
					"event_message": event.EventMessage,
					"entries":       event.Entries,
				},
			}
			bt.client.Publish(beatEvent)
			logp.Info("Beat Event sent")
		}
	}
}

func (bt *Crowdbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
