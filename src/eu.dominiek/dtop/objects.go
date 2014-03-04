package main

// These events are serialized and passed on to the clients.
type Event struct {
	Q string      // qualifier
	V interface{} // value (can be any struct)
}

type BasicInfo struct {
	Hostname         string
	SystemInfo       string
	DistributionInfo string
}

type DiskInfo struct {
	Name       string
	Type       string
	Size       string
	Used       string
	Available  string
	UsedPct    string
	MountPoint string
}

type CpuUsage struct {
	CpuId int
	Usage float64
}

type MemoryUsage struct {
	TotalKb     int
	FreeKb      int
	SharedKb    int
	BuffersKb   int
	CachedKb    int
	SwapTotalKb int
	SwapFreeKb  int
}

type LoadAverage struct {
	Avg1  float64
	Avg5  float64
	Avg15 float64
}

type ProcessInfo struct {
	Pid     int
	User    string
	Pri     int
	Ni      int
	Virt    int
	Res     int
	Shr     int
	S       string
	Cpu     float64
	Mem     float64
	Time    int64
	Command string
}

type User struct {
	Name string
}

type Users struct {
	All []User
}

func NewEvent(qualifier string, value interface{}) Event {
	event := new(Event)
	event.Q = qualifier
	event.V = value
	return *event
}

func NewBasicInfo(hostname string, systemInfo string, distributionInfo string) BasicInfo {
	basicInfo := new(BasicInfo)
	basicInfo.Hostname = hostname
	basicInfo.SystemInfo = systemInfo
	basicInfo.DistributionInfo = distributionInfo
	return *basicInfo
}

func NewMemoryUsage(totalKb int, freeKb int, sharedKb int, buffersKb int, cachedKb int, swapTotalKb int, swapFreeKb int) MemoryUsage {
	memoryUsage := new(MemoryUsage)
	memoryUsage.TotalKb = totalKb
	memoryUsage.FreeKb = freeKb
	memoryUsage.SharedKb = sharedKb
	memoryUsage.BuffersKb = buffersKb
	memoryUsage.CachedKb = cachedKb
	memoryUsage.SwapTotalKb = swapTotalKb
	memoryUsage.SwapFreeKb = swapFreeKb
	return *memoryUsage
}

func NewCpuUsage(cpuId int, usage float64) CpuUsage {
	cpuUsage := new(CpuUsage)
	cpuUsage.CpuId = cpuId
	cpuUsage.Usage = usage
	return *cpuUsage
}

func NewLoadAverage(avg1 float64, avg5 float64, avg15 float64) LoadAverage {
	loadAverage := new(LoadAverage)
	loadAverage.Avg1 = avg1
	loadAverage.Avg5 = avg5
	loadAverage.Avg15 = avg15
	return *loadAverage
}

func NewUser(name string) User {
	user := new(User)
	user.Name = name
	return *user
}

func NewUsers(all []User) Users {
	users := new(Users)
	users.All = all
	return *users
}

func NewProcessInfo(pid int, user string, pri int, ni int, virt int, res int, shr int, s string, cpu float64, mem float64, time int64, command string) ProcessInfo {
	processInfo := new(ProcessInfo)
	processInfo.Pid = pid
	processInfo.User = user
	processInfo.Pri = pri
	processInfo.Ni = ni
	processInfo.Virt = virt
	processInfo.Res = res
	processInfo.Shr = shr
	processInfo.S = s
	processInfo.Cpu = cpu
	processInfo.Mem = mem
	processInfo.Time = time
	processInfo.Command = command
	return *processInfo
}

func NewDiskInfo(name string, diskType string, size string, used string, available string, usedPct string, mountPoint string) DiskInfo {
	diskInfo := new(DiskInfo)
	diskInfo.Name = name
	diskInfo.Type = diskType
	diskInfo.Size = size
	diskInfo.Used = used
	diskInfo.Available = available
	diskInfo.UsedPct = usedPct
	diskInfo.MountPoint = mountPoint
	return *diskInfo
}
