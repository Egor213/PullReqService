package service

import (
	"context"

	"app/internal/repo"
	sd "app/internal/service/dto"
	smap "app/internal/service/mappers"
	errutils "app/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type StatsService struct {
	statsRepo repo.Stats
}

func NewStatsService(sRepo repo.Stats) *StatsService {
	return &StatsService{
		statsRepo: sRepo,
	}
}

func (s *StatsService) GetStats(ctx context.Context) (sd.GetStatsOutput, error) {
	statsPR, err := s.statsRepo.GetPRStats(ctx)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return sd.GetStatsOutput{}, err
	}

	statsReviewer, err := s.statsRepo.GetReviewerStats(ctx)
	if err != nil {
		log.Error(errutils.WrapPathErr(err))
		return sd.GetStatsOutput{}, err
	}
	return sd.GetStatsOutput{
		ByUsers: smap.ToReviewerStats(statsReviewer),
		ByPRs:   smap.ToPRStats(statsPR),
	}, nil
}
