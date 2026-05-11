//go:build windows

package server

func StartServer() error {
	return RunGoAsyncTCPServer()
}
