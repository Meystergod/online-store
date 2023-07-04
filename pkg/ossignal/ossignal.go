package ossignal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
)

type Signals []os.Signal

func (s Signals) String() string {
	ss := make([]string, 0, len(s))
	for _, sig := range s {
		ss = append(ss, sig.String())
	}

	return strings.Join(ss, ",")
}

type ErrSignal struct {
	Signal os.Signal
}

func (e ErrSignal) Error() string {
	return fmt.Sprintf("got error signal %s", e.Signal.String())
}

func DefaultOSSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT}
}

func DefaultSignalWaiter(ctx context.Context) error {
	return WaitSignal(ctx, DefaultOSSignals())
}

func WaitSignal(ctx context.Context, signals Signals) error {
	logger := zerolog.Ctx(ctx)

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, signals...)

	logger.Info().Stringer("signals", signals).Msg("wait os signals")
	select {
	case s := <-closeSignal:
		logger.Info().Msgf("got os signal: %s", s.String())
		return ErrSignal{Signal: s}
	case <-ctx.Done():
		return ctx.Err()
	}
}

func IsExitSignal(err error) bool {
	errSig := ErrSignal{}
	is := errors.As(err, &errSig)
	return is
}
