package handler

import (
	"hcc/cello/lib/logger"
	"os/exec"
	"strconv"
	"strings"
)

// ZSystem : Struct of ZSystem
type ZSystem struct {
	PoolName     string
	PoolCapacity string
	ZfsName      string
}

var zsysteminfo ZSystem

// CreateVolume : Creatte Volume
func CreateVolume(FileSystem string, ServerUUID string, VolType string, Size int) (bool, interface{}) {
	hostCheck()
	volcheck, err := QuotaCheck(ServerUUID)
	if !volcheck {
		logger.Logger.Println("CreateVolume : check Faild", err)
		return volcheck, err
	}
	createcheck, err := createzfs(FileSystem, ServerUUID, strings.ToUpper(VolType))
	if !createcheck {
		logger.Logger.Println("Create ZFS : Faild")
		return createcheck, err
	}
	setquota(ServerUUID, Size)

	return true, err

}

func createzfs(FileSystem string, ServerUUID string, VolType string) (bool, interface{}) {
	volname := FileSystem + VolType + "-vol-" + ServerUUID
	mountpath := "mountpoint=" + defaultdir + "/" + ServerUUID + "/" + FileSystem + "/" + VolType + "/"
	zsysteminfo.ZfsName = zsysteminfo.PoolName + "/" + volname
	cmd := exec.Command("zfs", "create", "-o", mountpath, zsysteminfo.ZfsName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

func setquota(ServerUUID string, Size int) (bool, interface{}) {
	qutoa := "quota=" + strconv.Itoa(Size) + "G"
	refquota := "refquota=" + strconv.Itoa(Size) + "G"

	cmd := exec.Command("zfs", "set", qutoa, refquota, zsysteminfo.ZfsName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(result, err)
	}
	zsysteminfo.ZfsName = ""
	return true, err
}
func hostCheck() {
	cmd := exec.Command("hostname")
	result, err := cmd.CombinedOutput()
	zsysteminfo.PoolName = strings.TrimSpace(string(result))
	if err != nil {
		logger.Logger.Println(result, err)
	}
}

//QuotaCheck : Zfs Available Quota check
func QuotaCheck(ServerUUID string) (bool, interface{}) {
	cmd := exec.Command("zfs", "get", "available", zsysteminfo.PoolName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	tmpstr := strings.Fields(string(result))
	var posofvalue int
	for i, words := range tmpstr {
		if words == "VALUE" {
			posofvalue = (len(tmpstr) / 2) + i
		}
	}
	zsysteminfo.PoolCapacity = tmpstr[posofvalue]
	return true, tmpstr[posofvalue]
}

// DeleteVolume :
// TODO : Implement delete volume
func DeleteVolume() {

}

// UpdateVolume :
// TODO : Implement update volume
func UpdateVolume() {

}
