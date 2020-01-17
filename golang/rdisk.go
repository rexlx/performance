package main

import (
	// "os/user"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
GIFTS FROM THE GOOD IDEA FAIRY

It may be nice to add a feature that, upon completion, crunches the
data and gives some generic info: min, max, mean, std... etc

considering adding colored threshholds when outputting to screen

DB support would be nice, considering mongo since i hate schemas.
likewise probably want to be able to send this data over api
*/

type Stat struct {
	Dev            string
	Time           int64
	Rsuccess       int
	Rmerged        int
	SectorRead     int
	Rtime          int
	Wcomplete      int
	Wmerged        int
	SectorWritten  int
	Wtime          int
	IOinProg       int
	IOtime         int
	WeightedTimeIO int
}

var (
	err            error
	name           string
	Dev            string
	Rsuccess       int
	Rmerged        int
	SectorRead     int
	Rtime          int
	Wcomplete      int
	Wmerged        int
	SectorWritten  int
	Wtime          int
	IOinProg       int
	IOtime         int
	WeightedTimeIO int
)

var entry Stat

// this is the help message we output when -h is passed or a bad flag is given
var help_msg string = `
usage: rdisk [-h] [-a] [-R REFRESH] [-r RUNTIME] 
            [-o OUTFILE]

This script records cpu statistics

optional arguments:
  -h          show this help message and exit
  -a          dont overwrite previous files
  -R          refresh rate, how long to wait between polls
  -r          runtime in seconds, or inf
  -d          send output to database (monogodb is supported)
`

// here we pasre the users args, return the args in a map
func parseArgs() (map[string]string, error) {
	// init the arg map
	argMap := make(map[string]string)
	// define our defaults
	argMap["append"] = "false"
	argMap["refresh"] = "5"
	argMap["runtime"] = "inf"
	argMap["db"] = "false"

	// call args
	args := os.Args
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

// this logs errors encountered in runtime
func diskLog(msg interface{}) error {
	file, err := os.OpenFile("rdisk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", file, ":", err)
		return nil
	}
	// defer file.Close()
	log.SetOutput(file)
	// test if input is string or error (or bad type)
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
		err := diskLog(e)
		if err != nil {
			fmt.Println("error logging to diskLog")
		}
		fmt.Println("encountered an error, check rcpu.log for details")
		os.Exit(1)
	}
}

// keep or whack the file as per cli args (or absence of)
func handleFile(header string, args map[string]string, filename string) error {
	// if we dont want to append (default in absense of the arg)
	if args["append"] == "false" {
		// the we want to try and remove it. if we cant, thats okay,
		// continue (file probably didnt exist)
		_file_ := filename + ".csv"
		err := os.Remove(_file_)
		if err != nil {
			fmt.Printf("couldnt remove %v: %v, continuing...\n", _file_, err)
		}
		// create the file now (or try, rather)
		fh, err := os.OpenFile(_file_, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("couldnt open %v: %v, continuing...\n", _file_, err)
		}
		// write the header and close the file
		fh.WriteString(header)
		fh.Close()
	}
	return nil
}

// writes our cpu util line
func writeFile(args map[string]string, line string, filename string) error {
	_file_ := filename + ".csv"
	file, err := os.OpenFile(_file_, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open file", file, ":", err)
		return nil
	}
	// write the cpu load line and close
	file.WriteString(line)
	file.Close()
	return nil
}

func send2db(args map[string]string, entry map[string]Stat) error {
	name, name_err := os.Hostname()
	if name_err != nil {
		return name_err
	}
	clientOptions := options.Client().ApplyURI(args["db"])
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	// err = client.Connect(ctx)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// fmt.Println("failed here")
		return err
	}
	collection := client.Database("disk_stats").Collection(name)
	res, res_err := collection.InsertOne(context.TODO(), entry)
	if res_err != nil {
		return res_err
	}
	log_err := diskLog(fmt.Sprintf("inserted %v", res))
	if log_err != nil {
		return log_err
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func poll_disks() map[string]Stat {
	// init our usage type
	usage := make(map[string]Stat)
	contents, err := ioutil.ReadFile("/proc/diskstats")
	check(err)
	// make array split by new line
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		// f is a slice of the current line we're on
		f := strings.Fields(line)
		if len(f) < 10 {
			continue
		}
		Rsuccess, err = strconv.Atoi(f[3])
		check(err)
		Rmerged, err = strconv.Atoi(f[4])
		check(err)
		SectorRead, err = strconv.Atoi(f[5])
		check(err)
		Rtime, err = strconv.Atoi(f[6])
		check(err)
		Wcomplete, err = strconv.Atoi(f[7])
		check(err)
		Wmerged, err = strconv.Atoi(f[8])
		check(err)
		SectorWritten, err = strconv.Atoi(f[9])
		check(err)
		Wtime, err = strconv.Atoi(f[10])
		check(err)
		IOinProg, err = strconv.Atoi(f[11])
		check(err)
		IOtime, err = strconv.Atoi(f[12])
		check(err)
		WeightedTimeIO, err = strconv.Atoi(f[13])
		check(err)
		entry = Stat{
			string(f[2]),
			0.0,
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
			WeightedTimeIO,
		}
		usage[string(f[2])] = entry
		// usage[string(line[2])] = entry
	}
	return usage
}

func main() {
	results := make(map[string]Stat)
	// for now we parse args in a sep func, it may make sense to do it
	// down here in main eventually
	args, parse_err := parseArgs()
	check(parse_err)
	// we call getCPUSample because we want to build a header and have
	// no idea how many cpus there are
	usage0 := poll_disks()
	keys := make([]string, 0, len(usage0))
	// hash tables are also unordered which can be a pain for
	// what we're doing, sort the keys so we iter over the same way
	// everytime
	for k := range usage0 {
		keys = append(keys, k)
		sort.Strings(keys)
	}
	// init some vars
	var duration float64
	var err error
	var entry Stat
	// build the header
	header := "utime," + strings.Join(keys[:], ",") + "\n"
	// overwrite or append the csv accordingly, pass it the header
	// we need to know how often to poll the file
	refresh, err := strconv.Atoi(args["refresh"])
	check(err)
	// go didnt have the kind of inf object i was hoping for, if
	// runtime is inf, its actually 1 billion seconds
	if args["runtime"] == "inf" {
		dur := &duration
		*dur = 1e9
	} else {
		// otherwise its what the user supplied
		dur := &duration
		*dur, err = strconv.ParseFloat(args["runtime"], 64)
		check(err)
	}
	// print the header line if we are providing stdout data
	// if args["silent"] == "false" {
	// 	for _, v := range keys {
	// 		fmt.Printf("%7v", v)
	// 	}
	// }
	var retry int
	// main loop starts here
	for i := time.Now(); time.Since(i) < time.Second*time.Duration(duration); {
		// probably could have init these up top
		var line string
		var lines []string
		var now string
		// get the unix time
		now = fmt.Sprintf("%v", time.Now().Unix())
		// unix time is our first field in each line
		lines = append(lines, now)
		// get our beginning poll
		usage0 = poll_disks()
		// wait for the desired refresh time
		time.Sleep(time.Duration(refresh) * time.Second)
		// get our ending poll
		usage1 := poll_disks()
		// we iter over the keys to ensure correct ordering
		for _, k := range keys {
			// this is how we calculate the usage
			Dev = usage1[k].Dev
			Rsuccess = usage1[k].Rsuccess - usage0[k].Rsuccess
			Rmerged = usage1[k].Rmerged - usage0[k].Rmerged
			SectorRead = usage1[k].SectorRead - usage0[k].SectorRead
			Rtime = usage1[k].Rtime - usage0[k].Rtime
			Wcomplete = usage1[k].Wcomplete - usage0[k].Wcomplete
			Wmerged = usage1[k].Wmerged - usage0[k].Wmerged
			SectorWritten = usage1[k].SectorWritten - usage0[k].SectorWritten
			Wtime = usage1[k].Wtime - usage0[k].Wtime
			IOinProg = usage1[k].IOinProg - usage0[k].IOinProg
			IOtime = usage1[k].IOtime - usage0[k].IOtime
			WeightedTimeIO = usage1[k].WeightedTimeIO - usage0[k].WeightedTimeIO
			entry = Stat{
				Dev,
				time.Now().Unix(),
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
				WeightedTimeIO,
			}
			results[Dev] = entry
		}
		// this tool in its current state ALWAYS creates a record
		if args["db"] != "false" {
			err := send2db(args, results)
			if err != nil {
				if retry >= 10 {
					exit_msg := "ran out of retries when connected to the DB, exiting..."
					fmt.Println(exit_msg)
					exit_err := diskLog(exit_msg)
					check(exit_err)
					os.Exit(1)
				}
				err := diskLog(err)
				check(err)
				time.Sleep(30 * time.Second)
				retry += 1
				info := fmt.Sprintf("encountered an error when connecting to the DB, waiting 30s and trying again (attempt %v/10)", retry)
				fmt.Println(info)
				e := diskLog(info)
				check(e)
				continue
			}
			retry = 0
		} else {
			for k, v := range results {
				handle_err := handleFile(header, args, k)
				check(handle_err)
				line = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
					v.Time, v.Rsuccess, v.Rmerged, v.SectorRead, v.Rtime, v.Wcomplete, v.Wmerged,
					v.SectorWritten, v.Wtime, v.IOinProg, v.IOtime, v.WeightedTimeIO)
				err := writeFile(args, line, k)
				check(err)
			}
			// write_err := writeFile(args, line)
			// check(write_err)
		}
	}
	// exit gracefully
	fmt.Println()
	os.Exit(0)
}
