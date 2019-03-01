package autoupdate

import (
	"github.com/aws/aws-sdk-go/aws"
)

type progressWriter struct {
    written int64
    writer *aws.WriteAtBuffer
    size int64
    progressCallback func(int)
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (n int, err error) {
	pw.written += int64(len(p))

	percentageDownloaded := int(float32(pw.written * 100) / float32(pw.size))

	if pw.progressCallback != nil {
		pw.progressCallback(percentageDownloaded)
	}

	return pw.writer.WriteAt(p, off)
}

