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
  -s          dont display statistics to screen
  -a          dont overwrite previous files
  -c          converts data to human readable
  -C          get the current usage in human readabale form
  -R          refresh rate, how long to wait between polls
  -r          runtime in seconds, or inf
  -o          outfile, file to write stats to
  -d          send output to database (monogodb is supported)
`

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
func memLog(msg interface{}) error {
	file, err := os.OpenFile("rmem.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
func check(e error) {
	if e != nil {
		err := memLog(e)
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
func convertToBytes(data string) (float64, error) {
	var new_int float64
	kib := math.Pow(1024, 1)
	mib := math.Pow(1024, 2)
	gib := math.Pow(1024, 3)
	tib := math.Pow(1024, 4)
	if strings.HasSuffix(data, "TB") {
		// data should be coming in as 666MB (for example), shave off the last 2 chars
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err)
		new_int = float64(converted_data) * tib
	} else if strings.HasSuffix(data, "GB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err)
		new_int = float64(converted_data) * gib
	} else if strings.HasSuffix(data, "MB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err)
		new_int = float64(converted_data) * mib
		// fmt.Println(converted_data)
	} else if strings.HasSuffix(data, "KB") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err)
		new_int = float64(converted_data) * kib
	} else if strings.HasSuffix(data, "B") {
		converted_data, err := strconv.Atoi(data[0 : len(data)-2])
		check(err)
		new_int = float64(converted_data)
	} else {
		e := errors.New("expected data to end with: B, KB, MB, GB, TB")
		err := memLog(e)
		check(err)
		return 0.0, e
	}
	return new_int, nil
}

func handleFile(header string, args map[string]string) error {
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

func getMem(args map[string]string) (string, error) {
	memMap := make(map[string]int)
	// memMap["el"] = 5
	fh, err := os.Open("/proc/meminfo")
	check(err)
	contents := bufio.NewScanner(fh)
	for contents.Scan() {
		line := strings.Fields(contents.Text())
		val, err := strconv.Atoi(line[1])
		check(err)
		memMap[line[0]] = val * 1024
	}
	fh.Close()
	_swap_ := memMap["SwapTotal:"] - memMap["SwapFree:"]
	_used_ := memMap["MemTotal:"] - memMap["MemFree:"] - memMap["Buffers:"] - memMap["Cached:"] - memMap["Slab:"]
	if args["current"] == "true" {
		if args["convert"] == "false" {
			message := fmt.Sprintf("free: %v\nused: %v\nswap: %v\n", memMap["MemFree:"], _used_, _swap_)
			return message, nil
		} else {
			free, err := convertBytes(memMap["MemFree:"])
			check(err)
			used, err := convertBytes(_used_)
			check(err)
			swap, err := convertBytes(_swap_)
			check(err)
			message := fmt.Sprintf("free: %v\nused: %v\nswap: %v\n", free, used, swap)
			return message, nil

		}
	}
	if args["convert"] == "false" {
		total := strconv.Itoa(memMap["MemTotal:"])
		free := strconv.Itoa(memMap["MemFree:"])
		buff := strconv.Itoa(memMap["Buffers:"])
		cache := strconv.Itoa(memMap["Cached:"])
		slab := strconv.Itoa(memMap["Slab:"])
		used := strconv.Itoa(_used_)
		swap := strconv.Itoa(_swap_)
		message := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n", time.Now().Unix(), total, free, buff, cache, slab, used, swap)
		return message, nil
	} else {
		total, err := convertBytes(memMap["MemTotal:"])
		check(err)
		free, err := convertBytes(memMap["MemFree:"])
		check(err)
		buff, err := convertBytes(memMap["Buffers:"])
		check(err)
		cache, err := convertBytes(memMap["Cached:"])
		check(err)
		slab, err := convertBytes(memMap["Slab:"])
		check(err)
		used, err := convertBytes(_used_)
		check(err)
		swap, err := convertBytes(_swap_)
		check(err)
		message := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n", time.Now().Unix(), total, free, buff, cache, slab, used, swap)
		return message, nil
	}
}

func metrics(quit chan bool) {
	for {
		select {
		case <-quit:
		default:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			util := fmt.Sprintf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\nHeap objects = %v\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC, m.Mallocs)
			err := memLog(util)
			check(err)
			time.Sleep(time.Second * 60)
		}
	}
}

func send2db(args map[string]string, line string) error {
	type Stat struct {
		Time  string
		Total string
		free  string
		Buff  string
		Cache string
		Slab  string
		Used  string
		Swap  string
	}
	stats := strings.Split(line, ",")
	entry := Stat{
		stats[0],
		stats[1],
		stats[2],
		stats[3],
		stats[4],
		stats[5],
		stats[6],
		stats[7],
	}
	name, name_err := os.Hostname()
	check(name_err)
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
	res, res_err := collection.InsertOne(context.TODO(), entry)
	if res_err != nil {
		return res_err
	}
	log_err := memLog(fmt.Sprintf("inserted %v", res))
	if log_err != nil {
		return log_err
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func main() {
	args, parse_err := parseArgs()
	check(parse_err)
	if args["current"] == "true" {
		mem, err := getMem(args)
		check(err)
		fmt.Println(mem)
		os.Exit(0)
	}
	quit := make(chan bool)
	go metrics(quit)
	var duration float64
	var err error
	header := "utime,total,used,free,buff,cache,slab,swap\n"
	handle_err := handleFile(header, args)
	check(handle_err)
	refresh, err := strconv.Atoi(args["refresh"])
	check(err)
	if args["runtime"] == "inf" {
		dur := &duration
		// *dur = math.Inf(1)
		*dur = 1e9
	} else {
		dur := &duration
		*dur, err = strconv.ParseFloat(args["runtime"], 64)
		check(err)
	}
	var retry int
	for i := time.Now(); time.Since(i) < time.Second*time.Duration(duration); {
		mem, err := getMem(args)
		check(err)
		if args["db"] != "false" {
			err := send2db(args, mem)
			if err != nil {
				if retry >= 10 {
					exit_msg := "ran out of retries when connected to the DB, exiting..."
					fmt.Println(exit_msg)
					exit_err := memLog(exit_msg)
					check(exit_err)
					os.Exit(1)
				}
				err := memLog(err)
				check(err)
				time.Sleep(30 * time.Second)
				retry += 1
				info := fmt.Sprintf("encountered an error when connecting to the DB, waiting 30s and trying again (attempt %v/10)", retry)
				fmt.Println(info)
				e := memLog(info)
				check(e)
				continue
			}
			retry = 0
		} else {
			mem += "\n"
			write_err := writeFile(args, mem)
			check(write_err)
		}
		time.Sleep(time.Duration(refresh) * time.Second)
	}
	quit <- true
	os.Exit(0)
}
