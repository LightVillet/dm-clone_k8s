package main

import (
	"os"
	"os/exec"
	"fmt"

	vph "k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
)

const (
	SectorSize = 512
)

func BlockDeviceSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func makeCloneTable(metaFile, destFile, sourceFile string) (string, error) {
	sectorCounts, err := BlockDeviceSize(sourceFile)
	if err != nil {
		return "", err
	}
	blkUtil := vph.NewBlockVolumePathHandler()
	metaDevice, err := blkUtil.AttachFileDevice(metaFile)
	if err != nil {
		return "", err
	}
	destDevice, err := blkUtil.AttachFileDevice(destFile)
	if err != nil {
		return "", err
	}
	sourceDevice, err := blkUtil.AttachFileDevice(sourceFile)
	if err != nil {
		return "", err
	}
	size := sectorCounts / SectorSize
	table := fmt.Sprintf("0 %d clone %s %s %s %d",
		size,
		metaDevice,
		destDevice,
		sourceDevice,
		size)
		return table, nil
}

func dmclone(name, metaFile, destFile, sourceFile string) error {
	table, err := makeCloneTable(metaFile, destFile, sourceFile)
	if err != nil {
		return err
	}
	_, err = exec.Command("dmsetup", "create", name, "--table", table).CombinedOutput()
	return nil
}

func main() {
	dmclone("dmclone", "/home/lightvillet/based/dm-clone_k8s/meta", "/home/lightvillet/based/dm-clone_k8s/dest", "/home/lightvillet/based/dm-clone_k8s/source")
}
