//go:build windows

package git

import "log"

func DeleteBranchesNotOnRemote(dryRun bool) (int, error) {
	log.Fatalln("DeleteBranchesNotOnRemote: windows not yet supported")
	return -1, nil
}
