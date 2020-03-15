package worker

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func (w *Worker) runTarget(job *Job) bool {
	job.Execute.StartTime = time.Now()

	cmd := exec.Command(filepath.Join(job.workingDir, job.ID+".exe"))
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		job.Execute.errChan <- "error piping to stdin from cmd: " + err.Error()
		w.log.Error("error piping stdin from cmd", zap.String("jobID", job.ID), zap.Error(err))
		return false
	}

	var output bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdout = &output

	if err := cmd.Start(); err != nil {
		job.Execute.errChan <- "error starting command: " + err.Error()
		w.log.Error("error starting command", zap.String("jobID", job.ID), zap.Error(err))
		return false
	}

	for _, line := range job.InputOutput.Input {
		if _, err := cmdStdin.Write([]byte(line + "\r\n\n")); err != nil {
			job.Execute.errChan <- "error writing to stdin: " + err.Error()
			w.log.Info("error writing to stdin", zap.String("jobID", job.ID), zap.Error(err))
		}
	}

	if err := cmdStdin.Close(); err != nil {
		job.Execute.errChan <- "error closing stdin pipe: " + err.Error()
		w.log.Error("error closing cmd stdin pipe", zap.String("jobID", job.ID), zap.Error(err))
	}

	// TODO: implement execution timeout with cmd.Process().Kill()
	if err := cmd.Wait(); err != nil {
		job.Execute.errChan <- err.Error()
		w.log.Error("error executing command", zap.String("jobID", job.ID), zap.Error(err))
	}

	scan := bufio.NewScanner(&output)
	for i := 0; scan.Scan(); i++ {
		str := scan.Text()
		if str != job.InputOutput.ExpectedOutput[i] {
			w.log.Debug("output mismatch", zap.String("jobID", job.ID),
				zap.String("expected", job.InputOutput.ExpectedOutput[i]),
				zap.String("received", str))
		}
		job.Execute.outChan <- str
	}

	job.Execute.EndTime = time.Now()
	job.Execute.Took = job.Execute.EndTime.Sub(job.Execute.StartTime)
	w.log.Info("executed target", zap.String("jobID", job.ID), zap.Duration("took", job.Execute.Took))

	return true
}
