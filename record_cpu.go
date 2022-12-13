package performance

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Usage struct {
	Total int
	Idle  int
}

type CpuUsage struct {
	Name  string
	Time  time.Time
	Usage float64
}

func GetCpuValues(c chan []*CpuUsage, refresh int) {
	if refresh < 1 {
		log.Println("cant wait less than 1 second")
		refresh = 1
	}

	var keys []string
	now := time.Now()
	values := []*CpuUsage{}

	initialPoll, err := pollCpu()
	for k := range initialPoll {
		keys = append(keys, k)
		sort.Strings(keys)
	}

	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Duration(refresh) * time.Second)
	poll, err := pollCpu()
	if err != nil {
		log.Println(err)
	}
	for _, key := range keys {
		idle := poll[key].Idle - initialPoll[key].Idle
		total := poll[key].Total - initialPoll[key].Total
		values = append(values, &CpuUsage{
			Name:  key,
			Usage: 100 * (float64(total) - float64(idle)) / float64(total),
			Time:  now})
	}
	c <- values
}

func pollCpu() (map[string]*Usage, error) {
	usage := make(map[string]*Usage)
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return usage, err
	}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		if strings.Contains(fields[0], "cpu") {
			result := &Usage{}
			nFields := len(fields)
			for i := 1; i < nFields; i++ {
				if i == 4 {
					val, err := strconv.Atoi(fields[i])
					if err != nil {
						return usage, err
					}
					result.Total += val
					result.Idle += val
				} else {
					val, err := strconv.Atoi(fields[i])
					if err != nil {
						return usage, err
					}
					result.Total += val
				}
				usage[fields[0]] = result
			}
		}
	}
	return usage, nil
}
