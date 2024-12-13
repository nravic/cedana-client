package db

// Local implementation of DB using SQL

import (
	"context"
	dbsql "database/sql"
	"os"
	"path/filepath"

	"buf.build/gen/go/cedana/cedana/protocolbuffers/go/daemon"
	"github.com/cedana/cedana/internal/db/sql"
	_ "github.com/mattn/go-sqlite3"
	json "google.golang.org/protobuf/encoding/protojson"
)

const name = "cedana.db"

type LocalDB struct {
	queries *sql.Queries
}

func NewLocalDB(ctx context.Context) (*LocalDB, error) {
	db, err := dbsql.Open("sqlite3", filepath.Join(os.TempDir(), name))
	if err != nil {
		return nil, err
	}

	// create sqlite tables
	if _, err := db.ExecContext(ctx, sql.Ddl); err != nil {
		return nil, err
	}

	return &LocalDB{
		queries: sql.New(db),
	}, nil
}

/////////////
// Getters //
/////////////

func (db *LocalDB) GetJob(ctx context.Context, jid string) (*daemon.Job, error) {
	dbJob, err := db.queries.GetJob(ctx, jid)
	if err != nil {
		return nil, err
	}

	bytes := dbJob.State

	// unmarsal the bytes into a Job struct
	job := daemon.Job{}
	err = json.Unmarshal(bytes, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

/////////////
// Setters //
/////////////

func (db *LocalDB) PutJob(ctx context.Context, jid string, job *daemon.Job) error {
	// marshal the Job struct into bytes
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	if _, err := db.queries.GetJob(ctx, jid); err == nil {
		return db.queries.UpdateJob(ctx, sql.UpdateJobParams{
			Jid:   jid,
			State: bytes,
		})
	} else {
		_, err := db.queries.CreateJob(ctx, sql.CreateJobParams{
			Jid:   jid,
			State: bytes,
		})
		return err
	}
}

/////////////
// Listers //
/////////////

func (db *LocalDB) ListJobs(ctx context.Context, jids ...string) ([]*daemon.Job, error) {
	dbJobs, err := db.queries.ListJobs(ctx)

	jidSet := make(map[string]struct{})
	for _, jid := range jids {
		jidSet[jid] = struct{}{}
	}

	jobs := []*daemon.Job{}
	for _, dbJob := range dbJobs {
		if len(jids) > 0 {
			if _, ok := jidSet[dbJob.Jid]; !ok {
				continue
			}
		}

		// unmarsal the bytes into a Job struct
		job := daemon.Job{}
		err = json.Unmarshal(dbJob.State, &job)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

//////////////
// Deleters //
//////////////

func (db *LocalDB) DeleteJob(ctx context.Context, jid string) error {
	return db.queries.DeleteJob(ctx, jid)
}