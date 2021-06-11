package base

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"io/ioutil"
	"log"
	"net"
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

func PrintViewUrl(port int) {
	if addr, err := net.InterfaceAddrs(); err == nil {
		log.Println("可通过以下地址访问")
		log.Printf("http://127.0.0.1:%d/web\n", port)
		for _, a := range addr {
			if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					log.Printf("http://%s:%d/web\n", ipNet.IP.String(), port)
				}
			}
		}
	} else {
		log.Println(err)
	}
}
