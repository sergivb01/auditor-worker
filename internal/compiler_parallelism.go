package worker

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

func (w *Worker) compileSingleFile(file string, wg *sync.WaitGroup) ([]byte, error) {
	defer w.Track("compileSingleFile("+file+")", Start())
	cmd := exec.Command("g++", "-O0", "-c", file+".cpp", "-o", file+".o")
	w.log.Info("executing cmd", zap.String("cmd", cmd.Path), zap.Strings("args", cmd.Args))
	defer wg.Done()
	return cmd.CombinedOutput()
}

func (w *Worker) compileWithParallelism(job *Job) {
	defer w.Track("compileFiles(job *Job)", Start())

	job.Build.StartTime = time.Now()

	var args = []string{"-std=c++11", "-O0"}

	var wg sync.WaitGroup
	wg.Add(len(job.CompileFiles))

	for _, file := range job.CompileFiles {
		go func(file string) {
			b, err := w.compileSingleFile(file, &wg)
			if err != nil {
				job.Build.errChan <- err.Error()
				return
			}
			job.Build.outChan <- string(b)
		}(file)
		args = append(args, file+".o")
	}
	wg.Wait()

	cmd := exec.Command("g++", append(args, "-o", filepath.Join(job.workingDir, "target.exe"))...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		job.Build.errChan <- err.Error() + fmt.Sprintf(cmd.Path+" (%s)", cmd.Args)
	}
	job.Build.outChan <- string(b)

	job.Build.EndTime = time.Now()
	job.Build.Took = job.Build.EndTime.Sub(job.Build.StartTime)

	w.log.Info("built project", zap.String("jobID", job.ID), zap.Duration("took", job.Build.Took))
}
