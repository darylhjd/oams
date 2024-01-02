package intervention

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/logger"
	"github.com/darylhjd/oams/backend/pkg/azmail"
	"go.uber.org/zap"
)

const (
	Namespace = "intervention"
)

type Service struct {
	l  *zap.Logger
	db *database.DB

	mailer *azmail.Client
}

// New creates the intervention service.
func New(ctx context.Context) (*Service, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	db, err := database.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - could not connect to database: %w", Namespace, err)
	}

	mailer, err := azmail.NewClient(
		os.Getenv("AZURE_EMAIL_ENDPOINT"),
		os.Getenv("AZURE_EMAIL_ACCESS_KEY"),
		os.Getenv("AZURE_EMAIL_SENDER_ADDRESS"),
	)
	if err != nil {
		return nil, fmt.Errorf("%s - could not create mailer client: %w", Namespace, err)
	}

	return &Service{
		l, db, mailer,
	}, nil
}

func (s *Service) Run() error {
	s.l.Info(fmt.Sprintf("%s - intervention service invoked", Namespace), zap.Time("time", time.Now()))

	if err := s.doStuff(); err != nil {
		return err
	}

	s.l.Info(fmt.Sprintf("%s - intervention service completed", Namespace), zap.Time("time", time.Now()))
	return nil
}

func (s *Service) doStuff() error {
	return errors.New("not implemented")
}

// Stop the intervention service gracefully.
func (s *Service) Stop() error {
	return s.db.Close()
}
