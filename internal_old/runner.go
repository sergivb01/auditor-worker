package worker

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func (w *Worker) runTarget(job *Job) error {
	defer w.Track("runTarget(job *Job)", Start())

	job.Execute.StartTime = time.Now()

	cmd := exec.Command(filepath.Join(job.workingDir, "target.exe"))
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		w.log.Error("error piping stdin from cmd", zap.String("jobID", job.ID), zap.Error(err))
		return fmt.Errorf("error piping to stdin from cmd: %w", err)
	}

	var output bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdout = &output

	if err := cmd.Start(); err != nil {
		w.log.Error("error starting command", zap.String("jobID", job.ID), zap.Error(err))
		return fmt.Errorf("error starting command: %w", err)
	}

	for _, line := range job.InputOutput.Input {
		if _, err := cmdStdin.Write([]byte(line + "\r\n\n")); err != nil {
			w.log.Info("error writing to stdin", zap.String("jobID", job.ID), zap.Error(err))
			return fmt.Errorf("error writing to stdin: %w", err)
		}
	}

	if err := cmdStdin.Close(); err != nil {
		w.log.Error("error closing cmd stdin pipe", zap.String("jobID", job.ID), zap.Error(err))
		return fmt.Errorf("error closing stdin pipe: %w", err)
	}

	// TODO: implement execution timeout with cmd.Process().Kill()
	if err := cmd.Wait(); err != nil {
		w.log.Error("error executing command", zap.String("jobID", job.ID), zap.Error(err))
		return fmt.Errorf("error exectuing command: %w", err)
	}

	scan := bufio.NewScanner(&output)
	for i := 0; scan.Scan(); i++ {
		str := scan.Text()
		if str != job.InputOutput.ExpectedOutput[i] {
			// job.Execute.errChan <- fmt.Sprintf("[line %d] output mismatch, expected %q but received %q", i+1, job.InputOutput.ExpectedOutput[i], str)
			job.Execute.Errors = append(job.Execute.Errors, fmt.Sprintf("[line %d] output mismatch, expected %q but received %q", i+1, job.InputOutput.ExpectedOutput[i], str))
			w.log.Debug("output mismatch", zap.String("jobID", job.ID),
				zap.String("expected", job.InputOutput.ExpectedOutput[i]),
				zap.String("received", str),
				zap.Int("line", i+1))
		}
		// job.Execute.outChan <- str
		job.Execute.Output = append(job.Execute.Output, str)
	}

	job.Execute.EndTime = time.Now()
	job.Execute.Took = job.Execute.EndTime.Sub(job.Execute.StartTime)
	w.log.Info("executed target", zap.String("jobID", job.ID), zap.Duration("took", job.Execute.Took))

	return nil
}
