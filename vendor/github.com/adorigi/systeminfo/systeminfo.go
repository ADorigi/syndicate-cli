package systeminfo

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
)

type AllInfo struct {
	CpuName         string  `json:"cpuname"`
	CpuArch         string  `json:"cpuarch"`
	OperatingSystem string  `json:"operatingsystem"`
	StorageUnit     string  `json:"storageunit"`
	DiskAvailable   uint64  `json:"diskavailable"`
	DiskUsed        uint64  `json:"DiskUsed"`
	DiskUsedPercent float64 `json:"DiskUsedPercent"`
	Hostname        string  `json:"Hostname"`
	LocalIPv4       string  `json:"LocalIPv4"`
	GlobalIP        string  `json:"GlobalIP"`
}

func NewAllInfo() AllInfo {
	return AllInfo{
		"",
		"",
		"",
		"GB",
		0,
		0,
		0.0,
		"",
		"",
		"",
	}
}

func (a *AllInfo) Initialize(cpuInfo cpu.InfoStat, dUsageInfo disk.UsageStat, hostInfo host.InfoStat) {

	a.CpuName = cpuInfo.ModelName
	a.CpuArch = runtime.GOARCH
	a.OperatingSystem = runtime.GOOS
	a.DiskAvailable = (uint64)(dUsageInfo.Free / 1000000000)
	a.DiskUsed = (uint64)(dUsageInfo.Used / 1000000000)
	a.DiskUsedPercent = dUsageInfo.UsedPercent
	a.Hostname = hostInfo.Hostname
	a.LocalIPv4 = GetLocalIPv4()
	a.GlobalIP = GetGlobalIP()
}

func GetLocalIPv4() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Could not gather IP")
		return "LocalIPv4 not found"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func GetGlobalIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		fmt.Println("GetGlobalIPv4 - Error in fetching public IP -- ", err)
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GetGlobalIPv4 - Error in reading response body -- ", err)
	}
	// fmt.Println(string(respData))
	return string(respData)
}

func CollectStats() AllInfo {

	allinfo := NewAllInfo()
	info, err := cpu.Info()
	if err != nil {
		fmt.Println("Error getting cpu info - ", err)
	}
	diskUsageInfo, err := disk.Usage("/")
	if err != nil {
		fmt.Println("Error getting disk info - ", err)
	}
	hostInfo, err := host.Info()
	if err != nil {
		fmt.Println("Error getting host info - ", err)
	}

	allinfo.Initialize(info[0], *diskUsageInfo, *hostInfo)

	return allinfo

}
