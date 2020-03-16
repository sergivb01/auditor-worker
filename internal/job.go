package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/sergivb01/acmecopy/api"
)

type JobTiming struct {
	StartTime time.Time
	EndTime   time.Time
	Took      time.Duration

	Output []string
	Errors []string

	// errChan chan string
	// outChan chan string
	//
	// sync.Mutex
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
		ID: id,
		workingDir: dir,
		Execute: JobTiming{
			// errChan: make(chan string),
			// outChan: make(chan string),
		},
		Build: JobTiming{
			// errChan: make(chan string),
			// outChan: make(chan string),
		},
		InputOutput: InputOutput{
			Input:          req.Input,
			ExpectedOutput: req.ExpectedOutput,
		},
	}, nil
}
//
// func (job *Job) listenForErrors(w *Worker) {
// 	for {
// 		select {
// 		case err := <-job.Build.errChan:
// 			job.Build.Lock()
// 			job.Build.Errors = append(job.Build.Errors, err)
// 			job.Build.Unlock()
// 			w.log.Info("received error from compiler", zap.String("jobID", job.ID), zap.String("error", err))
// 			break
//
// 		case err := <-job.Execute.errChan:
// 			job.Execute.Lock()
// 			job.Execute.Errors = append(job.Execute.Errors, err)
// 			job.Execute.Unlock()
// 			w.log.Info("received error from exec", zap.String("jobID", job.ID), zap.String("error", err))
// 			break
//
// 		}
// 	}
// }
//
// func (job *Job) listenForOutputs(w *Worker) {
// 	for {
// 		select {
// 		case buildOut := <-job.Build.outChan:
// 			if strings.TrimSpace(buildOut) != "" {
// 				job.Build.Lock()
// 				job.Build.Output = append(job.Build.Output, buildOut)
// 				job.Build.Unlock()
// 				w.log.Info("received output from compiler", zap.String("jobID", job.ID), zap.String("out", buildOut))
// 			}
// 			break
//
// 		case execOut := <-job.Execute.outChan:
// 			if strings.TrimSpace(execOut) != "" {
// 				job.Execute.Lock()
// 				job.Execute.Output = append(job.Execute.Output, execOut)
// 				job.Execute.Unlock()
// 				w.log.Info("received output from exec", zap.String("jobID", job.ID), zap.String("out", execOut))
// 			}
// 			break
// 		}
// 	}
// }
