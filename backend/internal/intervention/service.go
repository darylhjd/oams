package intervention

import (
	"context"
	"fmt"
	"time"

	"github.com/darylhjd/azmail"
	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/internal/logger"
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
		env.GetAzureEmailEndpoint(),
		env.GetAzureEmailAccessKey(),
		env.GetAzureEmailSenderAddress(),
	)
	if err != nil {
		return nil, fmt.Errorf("%s - could not create mailer client: %w", Namespace, err)
	}

	return &Service{
		l, db, mailer,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	s.l.Info(fmt.Sprintf("%s - intervention service invoked", Namespace), zap.Time("time", time.Now()))

	facts, rules, err := s.db.Intervention(ctx)
	if err != nil {
		return err
	}

	s.l.Info(
		fmt.Sprintf("%s - retrieved data from database", Namespace),
		zap.Int("num_facts", len(facts)),
		zap.Int("num_rules", len(rules)),
	)

	factGroups, ruleGroups := s.groupFacts(facts), s.groupRules(rules)
	users, ruleCreators, err := s.performChecks(factGroups, ruleGroups)
	if err != nil {
		return err
	}

	mails, err := s.generateNotificationMails(users, ruleCreators)
	if err != nil {
		return err
	}

	s.l.Info(fmt.Sprintf("%s - sending notification mails", Namespace), zap.Int("num_mails", len(mails)))

	if err = s.mailer.SendMails(mails...); err != nil {
		s.l.Warn(
			fmt.Sprintf("%s - some errors were encountered while sending notification mails", Namespace),
			zap.Error(err),
		)
	}

	s.l.Info(fmt.Sprintf("%s - intervention service completed", Namespace), zap.Time("time", time.Now()))
	return nil
}

// Stop the intervention service gracefully.
func (s *Service) Stop() error {
	return s.db.Close()
}
