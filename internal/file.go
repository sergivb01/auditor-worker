package worker

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/sergivb01/acmecopy/api"
	"go.uber.org/atomic"
)

func (job *Job) parseAndWriteFiles(req *api.CompileRequest) error {
	var (
		wg       sync.WaitGroup
		writeErr error
		success  = atomic.NewBool(true)
	)

	for _, file := range req.Files {
		if !success.Load() {
			break
		}

		go func(wg *sync.WaitGroup, file *api.File) {
			if filepath.Ext(file.FileName) == ".cpp" {
				job.CompileFiles = append(job.CompileFiles, filepath.Join(job.workingDir, removeExtension(file.FileName)))
			}
			job.UploadFiles = append(job.UploadFiles, filepath.Join(job.workingDir, file.FileName))

			if err := ioutil.WriteFile(filepath.Join(job.workingDir, file.FileName), file.Content, 744); err != nil {
				success.Store(false)
				writeErr = err
				return
			}
			wg.Done()
		}(&wg, file)
		wg.Add(1)
	}
	wg.Wait()

	return writeErr
}
