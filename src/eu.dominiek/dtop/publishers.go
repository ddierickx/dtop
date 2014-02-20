package main

import (
    "fmt"
    "time"
    "errors"
    "io/ioutil"
    "strings"
    "os"
    "os/exec"
    "regexp"
    "strconv"
)

// Default delay for publisher data refresh
const DELAY = 1 * time.Second

// Function template for event publishers.
type EventPublisher func(events chan Event)

// Process info publisher function which parses `ps` output into a ProcessInfo array and broadcasts it on the event channel.
func processinfo(events chan Event) {
   for {
        output, err := capture_stdout("ps", "auxh")

        if err != nil {
            panic(fmt.Sprintf("Error obtaining process info: %s", err))
        }

        var processInfos []ProcessInfo

        // iterate over each line, parsing is a bit annoying, basically space marks a new column
        // but the last column 'command' (10) can contain spaces so cut that one out.
        for _, line := range output {
            var values []string
            prevSpace := false
            value := ""

            // tokenize column values
            for i, char := range line {
                // command can have spaces in it but is the last column.
                if len(values) == 10 {
                    values = append(values, line[i:])
                    break
                }

                if char == ' ' && !prevSpace {
                    values = append(values, strings.TrimSpace(value))
                    value = ""
                    prevSpace = true
                } else if char != ' ' && prevSpace {
                    value = value + string(char)
                    prevSpace = false
                } else {
                    value = value + string(char)
                }
            }

            // Example header: USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND */
            pid := 0
            user := ""
            cpu := 0.0
            mem := 0.0
            command := ""

            // TODO: parse missing values compared to htop.
            for column, value := range values {
                switch column {
                    case 0:
                        user = value
                    case 1:
                        pid = atoi(value)
                    case 2:
                        cpu = atof(value)
                    case 3:
                        mem = atof(value)
                    case 10:
                        command = value
                }
            }

            processInfo := NewProcessInfo(pid, user, 0, 0, 0, 0, 0, "", cpu, mem, 0, command)
            processInfos = append(processInfos, processInfo)
        }

        events <- NewEvent("sys.processes", processInfos)

        time.Sleep(DELAY)
    }    
}

// Parse '/proc/stat' to find cpu usage per core.
func cpuinfo(events chan Event) {
    singlespace := regexp.MustCompile("\\ +")
    prevTotal := make(map[int]int)
    prevIdle := make(map[int]int)
    currentUsage := make(map[int]float64)

    for {
        data, err := ioutil.ReadFile("/proc/stat")

        if err == nil {
            lines := strings.Split(string(data), "\n")
            i := -1

            for _, line := range lines {
                if strings.HasPrefix(line, "cpu ") {
                    // skip summed total
                    continue
                }

                if strings.HasPrefix(line, "cpu") {
                    i = i + 1
                    columns := strings.Split(singlespace.ReplaceAllString(line, " "), " ")
                    userMode := atoi(columns[1])
                    userModeNice := atoi(columns[2])
                    system := atoi(columns[3])
                    idle := atoi(columns[4])
                    iowait := atoi(columns[5])
                    irq := atoi(columns[6])
                    softirq := atoi(columns[7])
                    steal := atoi(columns[8])
                    guest := atoi(columns[9])

                    newTotal := userMode + userModeNice + system + iowait + irq + softirq + steal + guest
                    newIdle := idle

                    oldTotal := prevTotal[i]
                    oldIdle := prevIdle[i] 

                    deltaTotal :=  newTotal - oldTotal
                    deltaIdle := newIdle - oldIdle 

                    usage := ( (float64(deltaTotal) / float64(deltaTotal + deltaIdle) ) * 100.0 )
                    currentUsage[i] = usage

                    prevTotal[i] = newTotal
                    prevIdle[i] = newIdle
                } else {
                    break
                }
            }

            cpuUsages := make([]CpuUsage, len(currentUsage))

            for i, usage := range currentUsage {
                cpuUsages[i] = NewCpuUsage(i, usage)
            }

            events <- NewEvent("sys.cpu", cpuUsages)
        }

        time.Sleep(DELAY)
    }
}

// Parse '/proc/meminfo' to get memory info.
func memory(events chan Event) {
    for {
        data, err := ioutil.ReadFile("/proc/meminfo")

        if err == nil {
            lines := strings.Split(string(data), "\n")
            values := make(map[string]int)

            for _, line := range lines {
                columns := strings.Split(line, ":")

                if len(columns) >= 2 {
                    key := columns[0]
                    value := strings.TrimSpace(strings.Replace(columns[1], "kB", "", -1))
                    intValue, cerr := strconv.Atoi(value)

                    if cerr == nil {
                        values[key] = intValue
                    }
                }
            }

            totalKb := values["MemTotal"]
            freeKb := values["MemFree"]
            sharedKb := 0 // unused in modern kernels but kept here for reference.
            buffersKb := values["Buffers"]
            cachedKb := values["Cached"]
            swapTotalKb := values["SwapTotal"]
            swapFreeKb := values["SwapFree"]

            events <- NewEvent("sys.memory", NewMemoryUsage(totalKb, freeKb, sharedKb, buffersKb, cachedKb, swapTotalKb, swapFreeKb))
        }

        time.Sleep(DELAY)
    }
}

// Run 'w' command and parse output to get logged in users list.
func users(events chan Event) {
    for {
        output, err := capture_stdout("/usr/bin/w", "-h")

        if err != nil {
            panic(fmt.Sprintf("error obtaining user list: %s", err))
        }

        var all []User

        for _, line := range output {
            columns := strings.Split(line, " ")

            if columns[0] != "" {
                user := NewUser(columns[0])
                all = append(all, user)
            }
        }

        events <- NewEvent("sys.users", NewUsers(all))

        time.Sleep(DELAY)
    }    
}

// Parse '/proc/loadavg' to get system load averages.
func loadavg(events chan Event) {
    for {
        data, err := ioutil.ReadFile("/proc/loadavg")

        if err == nil {
            loadavg  := strings.Split(string(data), " ")
            avg1,  _ := strconv.ParseFloat(loadavg[0], 64)
            avg5,  _ := strconv.ParseFloat(loadavg[1], 64)
            avg15, _ := strconv.ParseFloat(loadavg[2], 64)

            events <- NewEvent("sys.loadavg", NewLoadAverage(avg1, avg5, avg15))
        } else {
            panic("LoadAvg: " + err.Error())
        }

        time.Sleep(DELAY)
    }
}

// Run 'uptime' and parse output to get system uptime in seconds.
func uptime(events chan Event) {
    for {
        data, err := ioutil.ReadFile("/proc/uptime")

        if err == nil {
            uptime := strings.Split(string(data), " ")[0]
            events <- NewEvent("sys.uptime", uptime)
        } else {
            panic("Uptime: " + err.Error())
        }

        time.Sleep(DELAY)
    }
}

// Use internal go os module to get the system hostname.
func basicinfo(events chan Event) {
    hostname, hn_err := os.Hostname()

    if hn_err != nil {
        panic(fmt.Sprintf("Unable to get hostname: '%s'", hn_err))
    }

    release, rl_err := capture_stdout("lsb_release", "-sd")

    if rl_err != nil {
        panic(fmt.Sprintf("Unable to run lsb_release: '%s'", rl_err))
    }

    uname, un_err := capture_stdout("uname", "-sr")

    if un_err != nil {
        panic(fmt.Sprintf("Unable to run uname: '%s'", un_err))
    }

    basicInfo := NewBasicInfo(hostname, uname[0], release[0])
    events <- NewEvent("sys.basics", basicInfo)
}

// Capture stdout when running the cmd with the given arguments. Output is split on newline.
func capture_stdout(cmd string, args string) ([]string, error) {
        command := exec.Command(cmd, args)
        stdout, err := command.StdoutPipe()

        if err != nil {
            return nil, errors.New(fmt.Sprintf("Error redirecting stdout: %s", err))
        }

        if err := command.Start(); err != nil {
            return nil, errors.New(fmt.Sprintf("Error running command '%s' : %s", cmd, err))
        }

        output := ""

        if b, err := ioutil.ReadAll(stdout); err == nil {
            output = string(b)
        } else {
            return nil, errors.New(fmt.Sprintf("Error reading output stream: '%s'", err))
        }

        if err := command.Wait(); err != nil {
            return nil, errors.New(fmt.Sprintf("Error running command '%s': %s", cmd, err))
        }

        return strings.Split(output, "\n"), nil
}

// Convert the input string to a float64 and panic when parsing fails.
func atof(a string) float64 {
    x, err := strconv.ParseFloat(a, 64)

    if  err != nil {
        panic("Unable to convert to int: '" + a + "'")
    } else {
        return x
    }
}

// Convert the input string to an int (32) and panic when parsing fails.
func atoi(a string) int {
    i, err := strconv.Atoi(a)

    if  err != nil {
        panic("Unable to convert to int: '" + a + "'")
    } else {
        return i
    }
}