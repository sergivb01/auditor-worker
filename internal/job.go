package worker

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergivb01/acmecopy/api"
	"go.uber.org/zap"
)

type JobTiming struct {
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Took      time.Duration `json:"took"`

	errChan chan string `json:"-"`
	outChan chan string `json:"-"`
	//
	Output []string `json:"output"`
	Errors []string `json:"errors"`
}

type InputOutput struct {
	Input          []string `json:"input"`
	ExpectedOutput []string `json:"expectedInput"`
}

type Job struct {
	ID string `json:"id"`

	UploadFiles  []string `json:"uploadFiles"`
	CompileFiles []string `json:"compileFiles"`

	InputOutput InputOutput `json:"inputOutput"`

	req        *api.CompileRequest `json:"compileRequest"`
	workingDir string              `json:"workingDir"`

	Build   JobTiming `json:"build"`
	Execute JobTiming `json:"execute"`
}

func newJob(req *api.CompileRequest, dir string) *Job {
	return &Job{
		ID:         "abcdef",
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
	}
}

func (job *Job) listenForOutputs(w *Worker) {
	for {
		select {
		case err := <-job.Build.errChan:
			job.Build.Errors = append(job.Build.Errors, err)
			w.log.Info("received error from compiler", zap.String("jobID", job.ID), zap.String("error", err))
			fmt.Printf("error compiling file: %s", err)
			break

		case buildOut := <-job.Build.outChan:
			if strings.TrimSpace(buildOut) != "" {
				job.Build.Errors = append(job.Build.Errors, buildOut)
				w.log.Info("received output from compiler", zap.String("jobID", job.ID), zap.String("out", buildOut))
			}
			break

		case err := <-job.Execute.errChan:
			job.Build.Errors = append(job.Build.Errors, err)
			w.log.Info("received error from exec", zap.String("jobID", job.ID), zap.String("error", err))
			break

		case buildOut := <-job.Execute.outChan:
			if strings.TrimSpace(buildOut) != "" {
				job.Build.Errors = append(job.Build.Errors, buildOut)
				w.log.Info("received output from exec", zap.String("jobID", job.ID), zap.String("out", buildOut))
			}
			break
		}
	}
}
