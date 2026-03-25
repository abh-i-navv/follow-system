package services

import (
	"context"
	"errors"
	"follow-system/internal/repository"

	"github.com/google/uuid"
)

type FollowService struct {
	Repo *repository.FollowRepo
}

func NewFollowService(repo *repository.FollowRepo) *FollowService {
	return &FollowService{Repo: repo}
}

func (s *FollowService) FollowUser(
	ctx context.Context,
	followerID,
	targetID uuid.UUID,
	idempotencyKey string,
) error {
	if followerID == targetID {
		return errors.New("cannot follow yourself")
	}

	return s.Repo.FollowUser(
		ctx,
		followerID,
		targetID,
		idempotencyKey,
	)
}

func (s *FollowService) UnfollowUser(
	ctx context.Context,
	followerID,
	targetID uuid.UUID,
) error {
	return s.Repo.UnfollowUser(
		ctx,
		followerID,
		targetID,
	)
}

func (s *FollowService) GetFollower(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	return s.Repo.GetFollower(ctx, userID)
}
