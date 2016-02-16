package shared

/*import (
	"io"
	"net/http"
)

func NewEncryptedWriterProxy(writer io.Writer, sessionToken []byte) *EncryptedWriterProxy {
	return &EncryptedWriterProxy{writer, sessionToken}
}

type EncryptedWriterProxy struct {
	writer       io.Writer
	sessionToken []byte
}

func (e *EncryptedWriterProxy) Write(p []byte) (n int, err error) {
	cipher, err := EncryptSymmetric(e.sessionToken, p)
	if err != nil {
		return 0, err
	}
	return e.writer.Write(cipher)
}

func (e *EncryptedWriterProxy) Flush() {
	if flusher, ok := e.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}*/
