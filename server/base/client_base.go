package base

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"io/ioutil"
	"strings"
)

func GetPath(parent string) ([]FilePath, error) {
	var filePaths []FilePath
	if parent == "" {
		if p, _, _, err := host.PlatformInformation(); err == nil {
			if strings.Contains(p, "Windows") {
				partitions, _ := disk.Partitions(false)
				for _, partition := range partitions {
					filePaths = append(filePaths, FilePath{
						Title:  partition.Mountpoint,
						Parent: "",
						Key:    partition.Mountpoint,
						Leaf:   false,
					})
				}
			} else {
				filePaths = append(filePaths, FilePath{
					Title:  "/",
					Parent: "",
					Key:    "/",
					Leaf:   false,
				})
			}
		} else {
			return nil, err
		}
	} else {
		files, err := ioutil.ReadDir(parent)
		if err == nil {
			for _, f := range files {
				if f.IsDir() {
					var key string
					if parent != "/" {
						key = fmt.Sprintf("%s/%s", parent, f.Name())
					} else {
						key = fmt.Sprintf("/%s", f.Name())
					}
					filePaths = append(filePaths, FilePath{
						Title:  f.Name(),
						Parent: parent,
						Key:    key,
						Leaf:   false,
					})
				}
			}
		} else {
			return nil, err
		}
	}
	return filePaths, nil
}

func GetSpace(path string) (uint64, error) {
	us, err := disk.Usage(path)
	if err == nil {
		return us.Free, nil
	} else {
		return 0, err
	}
}
