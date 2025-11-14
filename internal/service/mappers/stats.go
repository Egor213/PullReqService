package servmappers

import (
	rd "app/internal/repo/dto"
	sd "app/internal/service/dto"
)

func ToPRStats(repoStats []rd.PRStatsOutput) []sd.PRStatsDTO {
	out := make([]sd.PRStatsDTO, 0, len(repoStats))
	for _, pr := range repoStats {
		out = append(out, sd.PRStatsDTO{
			PullReqID:   pr.PullReqID,
			Assignments: pr.Assignments,
		})
	}
	return out
}

func ToReviewerStats(repoStats []rd.ReviewerStatsOutput) []sd.ReviewerStatsDTO {
	out := make([]sd.ReviewerStatsDTO, 0, len(repoStats))
	for _, r := range repoStats {
		out = append(out, sd.ReviewerStatsDTO{
			UserID:      r.UserID,
			Assignments: r.Assignments,
		})
	}
	return out
}
