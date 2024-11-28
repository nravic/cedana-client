// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sql

import (
	"context"
)

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (jid, data) VALUES (?, ?)
RETURNING jid, data
`

type CreateJobParams struct {
	Jid  string
	Data []byte
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob, arg.Jid, arg.Data)
	var i Job
	err := row.Scan(&i.Jid, &i.Data)
	return i, err
}

const deleteJob = `-- name: DeleteJob :exec
DELETE FROM jobs WHERE jid = ?
`

func (q *Queries) DeleteJob(ctx context.Context, jid string) error {
	_, err := q.db.ExecContext(ctx, deleteJob, jid)
	return err
}

const getJob = `-- name: GetJob :one
SELECT jid, data FROM jobs WHERE jid = ?
`

func (q *Queries) GetJob(ctx context.Context, jid string) (Job, error) {
	row := q.db.QueryRowContext(ctx, getJob, jid)
	var i Job
	err := row.Scan(&i.Jid, &i.Data)
	return i, err
}

const listJobs = `-- name: ListJobs :many
SELECT jid, data FROM jobs
`

func (q *Queries) ListJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.db.QueryContext(ctx, listJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(&i.Jid, &i.Data); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJob = `-- name: UpdateJob :exec
UPDATE jobs SET data = ? WHERE jid = ?
`

type UpdateJobParams struct {
	Data []byte
	Jid  string
}

func (q *Queries) UpdateJob(ctx context.Context, arg UpdateJobParams) error {
	_, err := q.db.ExecContext(ctx, updateJob, arg.Data, arg.Jid)
	return err
}
