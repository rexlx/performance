package performance

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// use named indexes instead of otherwise seemingly random integers
const (
	_ int = iota + 1
	device
	readSuccess
	readMerged
	sectorRead
	readTime
	writeComplete
	writeMerged
	sectorWritten
	writeTime
	ioInProg
	ioTime
	weightedTimeIo
)

// DiskStat type represents the utilization of a given partition
type DiskStat struct {
	Dev string
	Rsuccess,
	Rmerged,
	SectorRead,
	Rtime,
	Wcomplete,
	Wmerged,
	SectorWritten,
	Wtime,
	IOinProg,
	IOtime,
	WeightedTimeIO int
	Time time.Time
}

// GetDiskUsage polls storage device statistics for a given interval in seconds
func GetDiskUsage(c chan []*DiskStat, refresh int) {
	if refresh < 1 {
		log.Println("cant wait less than 1 second")
		refresh = 1
	}

	now := time.Now()
	var keys []string
	var diskStats []*DiskStat

	initialPoll := pollDisks()
	for _, k := range initialPoll {
		keys = append(keys, k.Dev)
		sort.Strings(keys)
	}

	time.Sleep(time.Duration(refresh) * time.Second)
	usagePoll := pollDisks()

	for _, k := range keys {
		diskStats = append(diskStats, &DiskStat{
			Dev:            k,
			Time:           now,
			Rsuccess:       usagePoll[k].Rsuccess - initialPoll[k].Rsuccess,
			Rmerged:        usagePoll[k].Rmerged - initialPoll[k].Rmerged,
			SectorRead:     usagePoll[k].SectorRead - initialPoll[k].SectorRead,
			Rtime:          usagePoll[k].Rtime - initialPoll[k].Rtime,
			Wcomplete:      usagePoll[k].Wcomplete - initialPoll[k].Wcomplete,
			Wmerged:        usagePoll[k].Wmerged - initialPoll[k].Wmerged,
			SectorWritten:  usagePoll[k].SectorWritten - initialPoll[k].SectorWritten,
			Wtime:          usagePoll[k].Wtime - initialPoll[k].Wtime,
			IOinProg:       usagePoll[k].IOinProg - initialPoll[k].IOinProg,
			IOtime:         usagePoll[k].IOinProg - initialPoll[k].IOinProg,
			WeightedTimeIO: usagePoll[k].WeightedTimeIO - initialPoll[k].WeightedTimeIO,
		})
	}
	c <- diskStats
}

// pollDisks reads and parses the /proc/diskstats file
func pollDisks() map[string]DiskStat {
	usage := make(map[string]DiskStat)
	contents, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		log.Println(err)
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		stat := DiskStat{
			Dev:            fields[device],
			Rsuccess:       valueToInteger(fields[readSuccess]),
			Rmerged:        valueToInteger(fields[readMerged]),
			SectorRead:     valueToInteger(fields[sectorRead]),
			Rtime:          valueToInteger(fields[readTime]),
			Wcomplete:      valueToInteger(fields[writeComplete]),
			Wmerged:        valueToInteger(fields[writeMerged]),
			SectorWritten:  valueToInteger(fields[sectorWritten]),
			Wtime:          valueToInteger(fields[writeTime]),
			IOinProg:       valueToInteger(fields[ioInProg]),
			IOtime:         valueToInteger(fields[ioTime]),
			WeightedTimeIO: valueToInteger(fields[weightedTimeIo]),
		}
		usage[fields[2]] = stat
	}
	return usage
}

func valueToInteger(s string) int {
	out, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
	}
	return out
}
