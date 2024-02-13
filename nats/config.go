package nats

import (
	"errors"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	Token string
}

func LoadConfigFromEnvironment() Config {
	return Config{
		Token: os.Getenv("NATS_TOKEN"),
	}
}

func (cfg Config) Connect() (*nats.Conn, error) {
	var options []nats.Option
	if cfg.Token != "" {
		options = append(options, nats.Token(cfg.Token))
	}
	return nats.Connect(nats.DefaultURL, options...)
}

func (cfg Config) SendCommand(subject, data string) error {
	nc, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer nc.Close()

	req := nats.NewMsg(subject)
	req.Data = []byte(data)

	response, err := nc.RequestMsg(req, 3*time.Second)
	if err != nil {
		return err
	}

	if string(response.Data) != "ok" {
		return errors.New(string(response.Data))
	}

	return nil
}

func (cfg Config) SendEvent(subject, data string) error {
	nc, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer nc.Close()

	req := nats.NewMsg(subject)
	req.Data = []byte(data)

	if err := nc.PublishMsg(req); err != nil {
		return err
	}

	return nil
}
