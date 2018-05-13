package utils

import (
	"os"
	osuser "os/user"
	"strconv"
)

func UpdateFilePerms(file string, mode string, owner string, group string) error {

	// update the mode
	if mode != "" {

		mode, _ := strconv.ParseInt(mode, 0, 32)
		mode32 := uint32(mode)
		err := os.Chmod(file, os.FileMode(mode32))
		if err != nil {
			return err
		}
	}

	if owner != "" && group == "" {
		userchange, err := osuser.Lookup(owner)
		if err != nil {
			return err
		}
		userid, _ := strconv.Atoi(userchange.Uid)
		groupid, _ := strconv.Atoi(userchange.Gid)
		err = os.Chown(file, userid, groupid)
		if err != nil {
			return err
		}
	} else {
		userchange, err := osuser.Lookup(owner)
		if err != nil {
			return err
		}
		groupchange, err2 := osuser.LookupGroup(group)
		if err2 != nil {
			return err
		}
		userid, _ := strconv.Atoi(userchange.Uid)
		groupid, _ := strconv.Atoi(groupchange.Gid)
		err = os.Chown(file, userid, groupid)
		if err != nil {
			return err
		}
	}
	return nil
}
