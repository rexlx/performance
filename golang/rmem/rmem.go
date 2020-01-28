package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this is the help message we output when -h is passed or a bad flag is given
var help_msg string = `
usage: rmem [-h] [-s] [-a] [-c] [-C] [-R REFRESH] [-r RUNTIME]
            [-o OUTFILE]

This script records memory statistics

optional arguments:
  -h          show this help message and exit
  -a          dont overwrite previous csv file
  -c          converts data to human readable
  -C          get the current usage in human readabale form
  -R          refresh rate, how long to wait between polls
  -r          runtime in seconds, or inf
  -o          outfile, file to write stats to
  -l          logfile, specify log path
  -d          send output to database (monogodb is supported)
`

type Stat struct {
	Time     int64
	RawStats map[string]int
	Free     map[string]int
}

var (
	duration float64
	err      error
	_swap_   int
	_used_   int
	message  string
)

// here we pasre the users args, return the args in a map
func parseArgs() (map[string]string, error) {
	// init the arg map
	argMap := make(map[string]string)
	// define our defaults
	argMap["silent"] = "false"
	argMap["convert"] = "false"
	argMap["current"] = "false"
	argMap["append"] = "false"
	argMap["refresh"] = "5"
	argMap["runtime"] = "inf"
	argMap["outfile"] = "memutil.csv"
	argMap["logfile"] = "rmem.log"
	argMap["db"] = "false"
	// call args
	args := os.Args
	// fmt.Println(args)
	// loop over the args and parse them (we skip arg 0 which is the scripts name)
	for i, a := range args[1:] {
		// if the arg doesnt start with "-" then its following an arg and needs to be stored, skip it
		if !strings.HasPrefix(a, "-") {
			continue
		} else if a == "-h" {
			fmt.Println(help_msg)
			os.Exit(0)
		} else if a == "-s" {
			argMap["silent"] = "true"
		} else if a == "-a" {
			argMap["append"] = "true"
		} else if a == "-c" {
			argMap["convert"] = "true"
		} else if a == "-C" {
			argMap["current"] = "true"
		} else if a == "-o" {
			// in the case the arg is a non bool type, then store the actual val following the arg
			argMap["outfile"] = args[i+2]
		} else if a == "-R" {
			argMap["refresh"] = args[i+2]
		} else if a == "-r" {
			argMap["runtime"] = args[i+2]
		} else if a == "-d" {
			argMap["db"] = args[i+2]
		} else if a == "-l" {
			argMap["logfile"] = args[i+2]
		} else {
			// otherwise we got an unexpected arg
			fmt.Println(help_msg)
			fmt.Printf("unexpected argument: %v\n", a)
			os.Exit(1)
		}
	}
	return argMap, nil
}

// this logs errors encountered in runtime and logs heap usage
func memLog(msg interface{}, logfile string) error {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", file, ":", err)
		return nil
	}
	// defer file.Close()
	log.SetOutput(file)
	switch v := msg.(type) {
	case string:
		// i didnt know how to handle v not being used :/
		fmt.Sprintf("%v", v)
		logmsg := &msg
		log.Println(*logmsg)
	case error:
		fmt.Sprintf("%v", v)
		logmsg := &msg
		log.Println(*logmsg)
	default:
		// logmsg := false
		fmt.Sprintf("%v", v)
		err := errors.New("received bad type, expected string or error")
		return err

	}
	file.Close()
	return nil
}

// error checker
func check(e error, logfile string) {
	if e != nil {
		err := memLog(e, logfile)
		if err != nil {
			fmt.Println("error logging to memLog")
		}
		fmt.Println("encountered an error, check rcpu.log for details")
		os.Exit(1)
	}
}

// convert bytes to human readable
func convertBytes(data int) (string, error) {
	// init some vars
	var symbol string
	var new_data float64
	// define our data sizes
	kib := math.Pow(1024, 1)
	mib := math.Pow(1024, 2)
	gib := math.Pow(1024, 3)
	tib := math.Pow(1024, 4)
	// convert the input into a float
	f := float64(data)
	// convert it
	if f > tib {
		symbol = "TB"
		new_data = f / tib
	} else if f > gib {
		symbol = "GB"
		new_data = f / gib
	} else if f > mib {
		symbol = "MB"
		new_data = f / mib
	} else if f > kib {
		symbol = "KB"
		new_data = f / kib
	} else if f >= 0 {
		symbol = "B"
		new_data = f
	} else {
		// otherwise we probably got a negative number or bad type
		fmt.Printf("encountered an error! expected an int >= 0, got %v of type %T\n", data, data)
		err := errors.New("received bad data, need int or float >= 0")
		return "none", err
	}
	// dont want that weird notation, convert it to something humans can read and return it
	formatted_data := strconv.FormatFloat(new_data, 'f', 2, 64)
	converted_data := formatted_data + symbol
	return converted_data, nil
}

// does the opposite of above
func convertToBytes(data string, logfile string) (float64, error) {
	var new_int float64
	kib := math.Pow(1024, 1)
	mib := math.Pow(1024, 2)
	gib := math.Pow(1024, 3)
	tib := math.Pow(1024, 4)
	if strings.HasSuffix(data, "TB") {
		// data should be coming in as 666MB (for example), shave off the last 2 chars
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err, logfile)
		new_int = float64(converted_data) * tib
	} else if strings.HasSuffix(data, "GB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err, logfile)
		new_int = float64(converted_data) * gib
	} else if strings.HasSuffix(data, "MB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err, logfile)
		new_int = float64(converted_data) * mib
		// fmt.Println(converted_data)
	} else if strings.HasSuffix(data, "KB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err, logfile)
		new_int = float64(converted_data) * kib
	} else if strings.HasSuffix(data, "B") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err, logfile)
		new_int = float64(converted_data)
	} else {
		e := errors.New("expected data to end with: B, KB, MB, GB, TB")
		err := memLog(e, logfile)
		check(err, logfile)
		return 0.0, e
	}
	return new_int, nil
}

func handleFile(header string, args map[string]string) error {
	if args["db"] != "false" {
		return nil
	}
	if args["append"] == "false" {
		err := os.Remove(args["outfile"])
		if err != nil {
			fmt.Printf("couldnt remove %v: %v, continuing...\n", args["outfile"], err)
		}
		fh, err := os.OpenFile(args["outfile"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("couldnt open %v: %v, continuing...\n", args["outfile"], err)
		}
		// defer fh.Close()
		fh.WriteString(header)
		fh.Close()
	}
	return nil
}

func writeFile(args map[string]string, line string) error {
	file, err := os.OpenFile(args["outfile"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln("Failed to open file", file, ":", err)
		return nil
	}
	// defer file.Close()
	file.WriteString(line)
	file.Close()
	return nil
}

func getMem(args map[string]string, logfile string) map[string]int {
	memMap := make(map[string]int)
	// memMap["el"] = 5
	fh, err := os.Open("/proc/meminfo")
	check(err, logfile)
	contents := bufio.NewScanner(fh)
	for contents.Scan() {
		line := strings.Fields(contents.Text())
		val, err := strconv.Atoi(line[1])
		check(err, logfile)
		memMap[line[0]] = val * 1024
	}
	fh.Close()
	return memMap
}

func metrics(quit chan bool, logfile string) {
	for {
		select {
		case <-quit:
		default:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			util := fmt.Sprintf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\nHeap objects = %v\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC, m.Mallocs)
			err := memLog(util, logfile)
			check(err, logfile)
			time.Sleep(time.Second * 60)
		}
	}
}

func send2db(args map[string]string, memMap map[string]int, free map[string]int, logfile string) error {
	entry := Stat{
		time.Now().Unix(),
		memMap,
		free,
	}
	name, err := os.Hostname()
	check(err, logfile)
	clientOptions := options.Client().ApplyURI(args["db"])
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	// err = client.Connect(ctx)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("failed here")
		return err
	}
	collection := client.Database("mem_stats").Collection(name)
	res, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		return err
	}
	err = memLog(fmt.Sprintf("inserted %v", res), args["logfile"])
	if err != nil {
		return err
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func main() {
	args, err := parseArgs()
	logfile := args["logfile"]
	free := make(map[string]int)
	check(err, logfile)
	if args["current"] == "true" {
		memMap := getMem(args, logfile)
		_swap_ = memMap["SwapTotal:"] - memMap["SwapFree:"]
		_used_ = memMap["MemTotal:"] - memMap["MemFree:"] - memMap["Buffers:"] - memMap["Cached:"] - memMap["Slab:"]
		if args["convert"] == "false" {
			message := fmt.Sprintf("free: %v\nused: %v\nswap: %v\n", memMap["MemFree:"], _used_, _swap_)
			fmt.Println(message)
			os.Exit(0)
		} else {
			free, err := convertBytes(memMap["MemFree:"])
			check(err, logfile)
			used, err := convertBytes(_used_)
			check(err, logfile)
			swap, err := convertBytes(_swap_)
			check(err, logfile)
			message := fmt.Sprintf("free: %v\nused: %v\nswap: %v\n", free, used, swap)
			fmt.Println(message)
			os.Exit(0)
		}
	}
	quit := make(chan bool)
	go metrics(quit, logfile)
	header := "utime,total,used,free,buff,cache,slab,swap\n"
	err = handleFile(header, args)
	check(err, logfile)
	refresh, err := strconv.Atoi(args["refresh"])
	check(err, logfile)
	if args["runtime"] == "inf" {
		dur := &duration
		// *dur = math.Inf(1)
		*dur = 1e9
	} else {
		dur := &duration
		*dur, err = strconv.ParseFloat(args["runtime"], 64)
		check(err, logfile)
	}
	var retry int
	for i := time.Now(); time.Since(i) < time.Second*time.Duration(duration); {
		memMap := getMem(args, logfile)
		check(err, logfile)
		if args["convert"] == "false" {
			free["swap"] = memMap["SwapTotal:"] - memMap["SwapFree:"]
			free["used"] = memMap["MemTotal:"] - memMap["MemFree:"] - memMap["Buffers:"] - memMap["Cached:"] - memMap["Slab:"]
			free["total"] = memMap["MemTotal:"]
			free["free"] = memMap["MemFree:"]
			free["buff"] = memMap["Buffers:"]
			free["cache"] = memMap["Cached:"]
			free["slab"] = memMap["Slab:"]
			// free["used"] := _used_
			// swap := _swap_
			message = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n",
				time.Now().Unix(),
				free["total"],
				free["free"],
				free["buff"],
				free["cache"],
				free["slab"],
				free["used"],
				free["swap"])
		} else {
			if args["db"] != "false" {
				fmt.Println("Cant send converted lines (str type) to the db, exiting...")
				os.Exit(1)
			}
			total, err := convertBytes(memMap["MemTotal:"])
			check(err, logfile)
			free, err := convertBytes(memMap["MemFree:"])
			check(err, logfile)
			buff, err := convertBytes(memMap["Buffers:"])
			check(err, logfile)
			cache, err := convertBytes(memMap["Cached:"])
			check(err, logfile)
			slab, err := convertBytes(memMap["Slab:"])
			check(err, logfile)
			used, err := convertBytes(_used_)
			check(err, logfile)
			swap, err := convertBytes(_swap_)
			check(err, logfile)
			message = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n", time.Now().Unix(), total, free, buff, cache, slab, used, swap)
		}
		if args["db"] != "false" {
			err := send2db(args, memMap, free, logfile)
			if err != nil {
				if retry >= 10 {
					exit_msg := "ran out of retries when connected to the DB, exiting..."
					fmt.Println(exit_msg)
					err := memLog(exit_msg, args["logfile"])
					check(err, logfile)
					os.Exit(1)
				}
				err := memLog(err, logfile)
				check(err, logfile)
				time.Sleep(30 * time.Second)
				retry += 1
				info := fmt.Sprintf("encountered an error when connecting to the DB, waiting 30s and trying again (attempt %v/10)", retry)
				fmt.Println(info)
				e := memLog(info, args["logfile"])
				check(e, logfile)
				continue
			}
			retry = 0
		} else {
			message += "\n"
			err := writeFile(args, message)
			check(err, logfile)
		}
		time.Sleep(time.Duration(refresh) * time.Second)
	}
	quit <- true
	os.Exit(0)
}
