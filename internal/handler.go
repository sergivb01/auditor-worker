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
	dir, err := w.generateDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not get tempdir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Println("could not remove tempdir", err.Error())
			w.log.Error("could not remove temp dir", zap.Error(err))
		}
	}()

	job := newJob(req, dir)

	for _, file := range req.Files {
		if filepath.Ext(file.FileName) != ".h" && filepath.Ext(file.FileName) != ".cpp" {
			continue
		}

		if filepath.Ext(file.FileName) == ".cpp" {
			job.CompileFiles = append(job.CompileFiles, filepath.Join(dir, removeExtension(file.FileName)))
		}
		job.UploadFiles = append(job.UploadFiles, filepath.Join(dir, file.FileName))

		if err := ioutil.WriteFile(filepath.Join(dir, file.FileName), file.Content, 744); err != nil {
			return &api.CompileResponse{}, err
		}
	}

	printJob(job)

	go job.listenForOutputs(w)
	job.Execute.outChan <- "fuckoff"

	w.compileFiles(job)
	w.runTarget(job)

	printJob(job)

	return &api.CompileResponse{
		Build: &api.Response{
			StartTime: job.Build.StartTime.Unix(),
			EndTime:   job.Build.EndTime.Unix(),
			Took:      int64(job.Build.Took),
		},
		Execute: &api.Response{
			StartTime: job.Execute.StartTime.Unix(),
			EndTime:   job.Execute.EndTime.Unix(),
			Took:      int64(job.Execute.Took),
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
