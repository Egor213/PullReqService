package service_test

import (
	repomocks "app/internal/mocks/repomock"
	rd "app/internal/repo/dto"
	"app/internal/service"
	sd "app/internal/service/dto"
	serverrs "app/internal/service/errors"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStatsService_GetStats(t *testing.T) {

	ctx := context.Background()

	prStats := []rd.PRStatsOutput{
		{PullReqID: "p1", Assignments: 4},
		{PullReqID: "p2", Assignments: 3},
	}
	revStats := []rd.ReviewerStatsOutput{
		{UserID: "u1", Assignments: 3},
		{UserID: "u2", Assignments: 5},
	}

	type mockBehavior func(s *repomocks.MockStats)

	testCases := []struct {
		name         string
		mockBehavior mockBehavior
		want         sd.GetStatsOutput
		wantErr      error
	}{
		{
			name: "OK",
			mockBehavior: func(s *repomocks.MockStats) {
				s.EXPECT().GetPRStats(ctx).Return(prStats, nil)
				s.EXPECT().GetReviewerStats(ctx).Return(revStats, nil)
			},
			want: sd.GetStatsOutput{
				ByPRs: []sd.PRStatsDTO{
					{PullReqID: "p1", Assignments: 4},
					{PullReqID: "p2", Assignments: 3},
				},
				ByUsers: []sd.ReviewerStatsDTO{
					{UserID: "u1", Assignments: 3},
					{UserID: "u2", Assignments: 5},
				},
			},
		},
		{
			name: "Cannot get PR stats",
			mockBehavior: func(s *repomocks.MockStats) {
				s.EXPECT().GetPRStats(ctx).Return(nil, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetPRStats,
		},
		{
			name: "Cannot get reviewer stats",
			mockBehavior: func(s *repomocks.MockStats) {
				s.EXPECT().GetPRStats(ctx).Return(prStats, nil)
				s.EXPECT().GetReviewerStats(ctx).Return(nil, errors.New("other error"))
			},
			wantErr: serverrs.ErrCannotGetReviewerStats,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mStats := repomocks.NewMockStats(ctrl)

			tc.mockBehavior(mStats)

			svc := service.NewStatsService(mStats)

			got, err := svc.GetStats(ctx)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
