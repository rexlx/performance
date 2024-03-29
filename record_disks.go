package performance

import (
	"log"
	"os"
	"sort"
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
	Dev            string    `json:"device"`
	Rsuccess       int       `json:"read_success"`
	Rmerged        int       `json:"read_merged"`
	SectorRead     int       `json:"sector_read"`
	Rtime          int       `json:"time_reading"`
	Wcomplete      int       `json:"write_complete"`
	Wmerged        int       `json:"write_merged"`
	SectorWritten  int       `json:"sector_written"`
	Wtime          int       `json:"time_writing"`
	IOinProg       int       `json:"current_io"`
	IOtime         int       `json:"time_io"`
	WeightedTimeIO int       `json:"weighted_time_io"`
	Time           time.Time `json:"time"`
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
			IOtime:         usagePoll[k].IOtime - initialPoll[k].IOtime,
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
			Rsuccess:       ValueToInteger(fields[readSuccess]),
			Rmerged:        ValueToInteger(fields[readMerged]),
			SectorRead:     ValueToInteger(fields[sectorRead]),
			Rtime:          ValueToInteger(fields[readTime]),
			Wcomplete:      ValueToInteger(fields[writeComplete]),
			Wmerged:        ValueToInteger(fields[writeMerged]),
			SectorWritten:  ValueToInteger(fields[sectorWritten]),
			Wtime:          ValueToInteger(fields[writeTime]),
			IOinProg:       ValueToInteger(fields[ioInProg]),
			IOtime:         ValueToInteger(fields[ioTime]),
			WeightedTimeIO: ValueToInteger(fields[weightedTimeIo]),
		}
		usage[fields[2]] = stat
	}
	return usage
}
