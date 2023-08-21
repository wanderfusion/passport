package waitlist

import (
	"github.com/rs/zerolog/log"

	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/repositories/waitlist"
)

type Service struct {
	WaitlistRepo *waitlist.Database
}

func New(conf *config.DatabaseConfig) *Service {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	waitlistRepo := waitlist.New(conf)

	svc := &Service{
		WaitlistRepo: waitlistRepo,
	}

	return svc
}

func (s *Service) AddToWaitlist(mail, name string) error {
	return s.WaitlistRepo.AddUser(mail, name)
}
