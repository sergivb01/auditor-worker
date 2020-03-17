package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/xid"
	"github.com/sergivb01/acmecopy/api"
)

type Job struct {
	ID string

	CompileFiles []string

	InputOutput InputOutput

	workingDir string

	Build   JobTiming
	Execute JobTiming
}

type JobTiming struct {
	StartTime time.Time
	EndTime   time.Time
	Took      time.Duration

	Output []string
	Errors []string
}

type InputOutput struct {
	Input          []string
	ExpectedOutput []string
}

func (w *Worker) newJob(req *api.CompileRequest) (*Job, error) {
	id := xid.New().String()

	dir := filepath.Join(w.buildsDir, "builds_"+id)
	if err := os.Mkdir(dir, 0700); err != nil {
		return nil, fmt.Errorf("could not create %q dir: %w", dir, err)
	}

	return &Job{
		ID:         id,
		workingDir: dir,
		InputOutput: InputOutput{
			Input:          req.Input,
			ExpectedOutput: req.ExpectedOutput,
		},
	}, nil
}
