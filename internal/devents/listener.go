package devents

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Listener struct {
	ctx     context.Context
	cli     *docker.Client
	logger  *zerolog.Logger
	counter *prometheus.CounterVec
}

func NewListener(ctx context.Context,
	cli *docker.Client, logger *zerolog.Logger, cfg *viper.Viper,
) *Listener {
	return &Listener{
		ctx:    ctx,
		cli:    cli,
		logger: logger,
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.GetString("metric.namespace"),
			Subsystem: cfg.GetString("metric.subsystem"),
			Name:      cfg.GetString("metric.name"),
			Help:      "Number of docker events.",
		}, []string{"type", "action", "status", "scope", "from", "name"}),
	}
}

func (l *Listener) Listen() error {
	eventsChan, errChan := l.cli.Events(l.ctx, types.EventsOptions{})

	for {
		select {
		case event := <-eventsChan:
			l.counter.WithLabelValues(event.Type, event.Action, event.Status,
				event.Scope, event.From, event.Actor.Attributes["name"]).Inc()
		case err := <-errChan:
			// if !errors.As(err, &io.EOF{}) { @TODO:
			// }
			return err
		}
	}
}
