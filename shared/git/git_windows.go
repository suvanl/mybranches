//go:build windows

package git

import "log"

func DeleteBranchesNotOnRemote() (int32, error) {
	log.Fatalln("DeleteBranchesNotOnRemote: windows not yet supported")
	return -1, nil
}
