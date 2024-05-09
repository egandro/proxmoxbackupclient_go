//go:build linux || darwin || freebsd || openbsd
// +build linux darwin freebsd openbsd

package pbs

func createVSSSnapshot(path string) string {
	return path
}

func VSSCleanup(){
}
