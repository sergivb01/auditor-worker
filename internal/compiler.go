package worker

import (
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

func compileSingleFile(file string, wg *sync.WaitGroup) ([]byte, error) {
	cmd := exec.Command("g++", "-c", file+".cpp", "-o", file+".o")
	defer wg.Done()
	return cmd.CombinedOutput()
}

func (w *Worker) compileFiles(job *Job) {
	job.Build.StartTime = time.Now()

	var args = []string{"-std=c++11"}

	var wg sync.WaitGroup
	wg.Add(len(job.CompileFiles))

	for _, file := range job.CompileFiles {
		go func(file string) {
			b, err := compileSingleFile(file, &wg)
			if err != nil {
				job.Build.errChan <- err.Error()
				return
			}
			job.Build.outChan <- string(b)
		}(file)
		args = append(args, file+".o")
	}
	wg.Wait()

	b, err := exec.Command("g++", append(args, "-o", filepath.Join(job.workingDir, job.ID+".exe"))...).CombinedOutput()
	if err != nil {
		job.Build.errChan <- err.Error()
	}
	job.Build.outChan <- string(b)

	job.Build.EndTime = time.Now()
	job.Build.Took = job.Build.EndTime.Sub(job.Build.StartTime)

	w.log.Info("built project", zap.String("jobID", job.ID), zap.Duration("took", job.Build.Took))
}
