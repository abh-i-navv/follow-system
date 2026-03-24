package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type FollowRepo struct {
	DB *sql.DB
}

func NewFollowRepo(db *sql.DB) *FollowRepo {
	return &FollowRepo{DB: db}
}

func (r *FollowRepo) FollowUser(
	ctx context.Context,
	followerID,
	targetID uuid.UUID,
	idempotencyKey string,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO follows
		(follower_id, following_id, status, idempotency_key)
		VALUES ($1, $2, 'active', $3)
		`,
		followerID,
		targetID,
		idempotencyKey,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE users
		 SET follower_count = follower_count+1
		 WHERE id = $1`,
		targetID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`
		UPDATE users
		SET following_count = following_count+1
		WHERE id = $1
	`,
		followerID,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}
