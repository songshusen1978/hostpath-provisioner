package hostpath

import (
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/kubernetes/pkg/volume/util/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (hpc *hostPathController) validateCreateVolumeCapacityRange(req *csi.CreateVolumeRequest) error {
	if req.CapacityRange == nil{
		return status.Error(codes.InvalidArgument, "must have capacity range")
	}
	if req.CapacityRange.RequiredBytes <= 0{
		return status.Error(codes.InvalidArgument, "capacity range required bytes <= 0")
	}
	return nil
}

func getPVStatsEx(volumePath string) (available int64, capacity int64, used int64, inodes int64, inodesFree int64, inodesUsed int64, err error) {
	return fs.Info(volumePath)
}

func diskdu(dir string)int64{
	info, err := os.Lstat(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return diskUsage(dir, info, 0)
}

func diskUsage(currPath string, info os.FileInfo, depth int) int64 {
	var size int64

	dir, err := os.Open(currPath)
	if err != nil {
		fmt.Println(err)
		return size
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			size += diskUsage(fmt.Sprintf("%s/%s", currPath, file.Name()), file, depth+1)
		} else {
			size += file.Size()
		}
	}

	return size
}
