package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/sergivb01/acmecopy/api"
	"go.uber.org/zap"
)

type JobTiming struct {
	StartTime time.Time
	EndTime   time.Time
	Took      time.Duration

	errChan chan string
	outChan chan string

	Output []string
	Errors []string
}

type InputOutput struct {
	Input          []string
	ExpectedOutput []string
}

type Job struct {
	ID string

	UploadFiles  []string
	CompileFiles []string

	InputOutput InputOutput

	req        *api.CompileRequest
	workingDir string

	Build   JobTiming
	Execute JobTiming

	sync.Mutex
}

func (w *Worker) newJob(req *api.CompileRequest) (*Job, error) {
	id := xid.New().String()

	dir := filepath.Join(w.buildsDir, "builds_"+id)
	if err := os.Mkdir(dir, 0700); err != nil {
		return nil, fmt.Errorf("could not create %q dir: %w", dir, err)
	}

	return &Job{
		ID:         id,
		req:        req,
		workingDir: dir,
		Execute: JobTiming{
			errChan: make(chan string),
			outChan: make(chan string),
		},
		Build: JobTiming{
			errChan: make(chan string),
			outChan: make(chan string),
		},
		InputOutput: InputOutput{
			Input: []string{
				"hola test123",
				"#",
			},
			ExpectedOutput: []string{
				"ENTRA TEXT ACABAT EN #:",
				"TEXT REVES:",
				"test123 hola ",
			},
		},
	}, nil
}

func (job *Job) listenForOutputs(w *Worker) {
	for {
		select {
		case err := <-job.Build.errChan:
			job.Build.Errors = append(job.Build.Errors, err)
			w.log.Info("received error from compiler", zap.String("jobID", job.ID), zap.String("error", err))
			break

		case buildOut := <-job.Build.outChan:
			if strings.TrimSpace(buildOut) != "" {
				job.Build.Output = append(job.Build.Output, buildOut)
				w.log.Info("received output from compiler", zap.String("jobID", job.ID), zap.String("out", buildOut))
			}
			break

		case err := <-job.Execute.errChan:
			job.Execute.Errors = append(job.Execute.Errors, err)
			w.log.Info("received error from exec", zap.String("jobID", job.ID), zap.String("error", err))
			break

		case execOut := <-job.Execute.outChan:
			if strings.TrimSpace(execOut) != "" {
				job.Execute.Output = append(job.Execute.Output, execOut)
				w.log.Info("received output from exec", zap.String("jobID", job.ID), zap.String("out", execOut))
			}
			break
		}
	}
}
