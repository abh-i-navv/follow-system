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

func (r *FollowRepo) UnfollowUser(
	ctx context.Context,
	followerID,
	targetID uuid.UUID,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`
		DELETE FROM follows
		WHERE follower_id = $1
		AND following_id = $2
	`,
		followerID,
		targetID,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return nil
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE users
		 SET follower_count = follower_count-1
		 WHERE id=$1`,
		targetID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE users
		 SET following_count = following_count - 1
		 WHERE id = $1`,
		followerID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FollowRepo) GetFollower(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {

	rows, err := r.DB.QueryContext(
		ctx,
		`SELECT follower_id
		FROM follows
		WHERE following_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var followers []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		followers = append(followers, id)
	}
	return followers, nil
}
