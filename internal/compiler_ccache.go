package worker

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func (w *Worker) compileWithCCache(job *Job) {
	defer w.Track("compileWithCCache(job *Job)", Start())

	job.Build.StartTime = time.Now()

	var args = []string{"g++", "-std=c++11", "-O0"}
	for _, file := range job.CompileFiles {
		args = append(args, file+".cpp")
	}

	cmd := exec.Command("ccache", append(args, "-o", filepath.Join(job.workingDir, "target.exe"))...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		job.Build.errChan <- err.Error() + fmt.Sprintf(cmd.Path+" (%s)", cmd.Args)
	}
	job.Build.outChan <- string(b)

	job.Build.EndTime = time.Now()
	job.Build.Took = job.Build.EndTime.Sub(job.Build.StartTime)

	w.log.Info("built project", zap.String("jobID", job.ID), zap.Duration("took", job.Build.Took))
}