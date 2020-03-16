package worker

import (
	"context"
	"fmt"
	"os"

	"github.com/sergivb01/acmecopy/api"
	"go.uber.org/zap"
)

func (w *Worker) CompileFiles(_ context.Context, req *api.CompileRequest) (*api.CompileResponse, error) {
	job, err := w.newJob(req)
	if err != nil {
		return nil, fmt.Errorf("could not create new job: %w", err)
	}

	if w.cfg.RemoveOld {
		defer func() {
			if err := os.RemoveAll(job.workingDir); err != nil {
				fmt.Println("could not remove tempdir", err.Error())
				w.log.Error("could not remove temp dir", zap.Error(err))
			}
		}()
	}

	if err := job.parseAndWriteFiles(req); err != nil {
		return nil, fmt.Errorf("error writing temp files: %w", err)
	}

	if err := w.compileWithCCache(job); err != nil {
		return nil, fmt.Errorf("error compiling: %w", err)
	}

	if err := w.runTarget(job); err != nil {
		return nil, fmt.Errorf("error running target: %w", err)
	}

	return &api.CompileResponse{
		Id: job.ID,
		Build: &api.Response{
			StartTime: job.Build.StartTime.Unix(),
			EndTime:   job.Build.EndTime.Unix(),
			Took:      int64(job.Build.Took),
			Errors:    job.Build.Errors,
			Log:       job.Build.Output,
		},
		Execute: &api.Response{
			StartTime: job.Execute.StartTime.Unix(),
			EndTime:   job.Execute.EndTime.Unix(),
			Took:      int64(job.Execute.Took),
			Errors:    job.Execute.Errors,
			Log:       job.Execute.Output,
		},
	}, nil
}
