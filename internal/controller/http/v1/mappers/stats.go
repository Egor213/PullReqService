package httpmappers

import (
	hd "app/internal/controller/http/v1/dto"
	sd "app/internal/service/dto"
)

func ToPRStats(repoStats []sd.PRStatsDTO) []hd.PRStatsDTO {
	out := make([]hd.PRStatsDTO, 0, len(repoStats))
	for _, pr := range repoStats {
		out = append(out, hd.PRStatsDTO{
			PullReqID:   pr.PullReqID,
			Assignments: pr.Assignments,
		})
	}
	return out
}

func ToReviewerStats(repoStats []sd.ReviewerStatsDTO) []hd.ReviewerStatsDTO {
	out := make([]hd.ReviewerStatsDTO, 0, len(repoStats))
	for _, r := range repoStats {
		out = append(out, hd.ReviewerStatsDTO{
			UserID:      r.UserID,
			Assignments: r.Assignments,
		})
	}
	return out
}
