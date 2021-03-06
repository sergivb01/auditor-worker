package worker

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

var ccachePath string

func init() {
	foundCcachePath, err := exec.LookPath("ccache")
	if err != nil {
		panic("couldn't find ccache in path!")
	}
	ccachePath = foundCcachePath
}

func (w *Worker) compileWithCCache(job *Job) error {
	defer w.Track("compileWithCCache(jobID="+job.ID+")", Start())

	job.Build.StartTime = time.Now()

	var args = []string{"g++", "-std=c++11", "-O0"}
	for _, file := range job.CompileFiles {
		args = append(args, file)
	}

	cmd := exec.Command(ccachePath, append(args, "-o", filepath.Join(job.workingDir, "target.exe"))...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		// job.Build.errChan <- err.Error() + fmt.Sprintf(cmd.Path+" (%s)", cmd.Args)
		return fmt.Errorf("error compiling with ccache: %w", err)
	}
	job.Build.Output = append(job.Build.Output, string(b))

	job.Build.EndTime = time.Now()
	job.Build.Took = job.Build.EndTime.Sub(job.Build.StartTime)

	return nil
}
