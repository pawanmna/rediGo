//go:build linux

package server

func StartServer() error {
	return RunAsyncTCPServer()
}
