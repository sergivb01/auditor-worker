package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

	// TODO: speed up this
	t := Start()
	for _, file := range req.Files {
		if filepath.Ext(file.FileName) != ".h" && filepath.Ext(file.FileName) != ".cpp" {
			continue
		}

		if filepath.Ext(file.FileName) == ".cpp" {
			job.CompileFiles = append(job.CompileFiles, filepath.Join(job.workingDir, removeExtension(file.FileName)))
		}
		job.UploadFiles = append(job.UploadFiles, filepath.Join(job.workingDir, file.FileName))

		if err := ioutil.WriteFile(filepath.Join(job.workingDir, file.FileName), file.Content, 744); err != nil {
			return &api.CompileResponse{}, err
		}
	}
	w.Track("parse and write files", t)

	go job.listenForOutputs(w)

	if w.cfg.CCacheEnabled {
		w.compileWithCCache(job)
	} else {
		w.compileWithParallelism(job)
	}
	w.runTarget(job)

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

func printJob(job *Job) {
	b, err := json.Marshal(job)
	if err != nil {
		fmt.Println("12334" + err.Error())
		return
	}
	fmt.Printf("%s\n\n", b)
}
