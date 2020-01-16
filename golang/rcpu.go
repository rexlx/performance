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

// create our usage type
type Usage struct {
	total, idle int
}

// this is the help message we output when -h is passed or a bad flag is given
var help_msg string = `
usage: rcpu [-h] [-s] [-a] [-c] [-C] [-R REFRESH] [-r RUNTIME]
            [-o OUTFILE]

This script records cpu statistics

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
	argMap["outfile"] = "cpuutil.csv"
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
func cpuLog(msg interface{}) error {
	file, err := os.OpenFile("rcpu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		err := cpuLog(e)
		if err != nil {
			fmt.Println("error logging to cpuLog")
		}
		fmt.Println("encountered an error, check rcpu.log for details")
		os.Exit(1)
	}
}

// keep or whack the file as per cli args (or absence of)
func handleFile(header string, args map[string]string) error {
	// if we dont want to append (default in absense of the arg)
	if args["append"] == "false" {
		// the we want to try and remove it. if we cant, thats okay,
		// continue (file probably didnt exist)
		err := os.Remove(args["outfile"])
		if err != nil {
			fmt.Printf("couldnt remove %v: %v, continuing...\n", args["outfile"], err)
		}
		// create the file now (or try, rather)
		fh, err := os.OpenFile(args["outfile"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("couldnt open %v: %v, continuing...\n", args["outfile"], err)
		}
		// write the header and close the file
		fh.WriteString(header)
		fh.Close()
	}
	return nil
}

// writes our cpu util line
func writeFile(args map[string]string, line string) error {
	file, err := os.OpenFile(args["outfile"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open file", file, ":", err)
		return nil
	}
	// write the cpu load line and close
	file.WriteString(line)
	file.Close()
	return nil
}

func send2db(args map[string]string, line string) error {
	type Stat struct {
		Time  string
		Stats []string
	}
	entry := Stat{strings.Split(line, ",")[0], strings.Split(line, ",")[1:]}
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
		fmt.Println("failed here")
		return err
	}
	collection := client.Database("cpu_stats").Collection(name)
	res, res_err := collection.InsertOne(context.TODO(), entry)
	if res_err != nil {
		return res_err
	}
	log_err := cpuLog(fmt.Sprintf("inserted %v", res))
	if log_err != nil {
		return log_err
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func poll_cpu() map[string]Usage {
	/*    user    nice   system  idle      iowait irq   softirq  steal  guest  guest_nice
	cpu  74608   2520   24433   1117073   6176   4054  0        0      0      0*/
	// init our usage type
	usage := make(map[string]Usage)
	// read the stat file
	contents, err := ioutil.ReadFile("/proc/stat")
	check(err)
	// make array split by new line
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		// f is a slice of the current line we're on
		f := strings.Fields(line)
		if len(f) < 1 {
			continue
		}
		// we only care about the lines with 'cpu' in them
		if strings.Contains(f[0], "cpu") {
			// init these for later
			total := 0
			idle := 0
			numFields := len(f)
			for i := 1; i < numFields; i++ {
				// these are the fields we consider idle
				if i == 4 {
					// if i == 4 || i == 5 {
					val, err := strconv.Atoi(f[i])
					check(err)
					// we want to add these to idle AND total or the
					// math dont work
					idle += val
					total += val
				} else {
					// otherwise its a line with cpu and a field we
					// consider part of total (but not idle)
					val, err := strconv.Atoi(f[i])
					check(err)
					total += val
				}
				// add that to our usage map as a Usage type with key:
				// 'cpu' or 'cpu0' (etc)
				usage[f[0]] = Usage{
					total,
					idle,
				}

			}
		}
	}
	return usage
}

func main() {
	// for now we parse args in a sep func, it may make sense to do it
	// down here in main eventually
	args, parse_err := parseArgs()
	check(parse_err)
	// we call getCPUSample because we want to build a header and have
	// no idea how many cpus there are
	usage0 := poll_cpu()
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
	var util string
	var idle int
	var total int
	var data []string
	// build the header
	header := "utime," + strings.Join(keys[:], ",") + "\n"
	// overwrite or append the csv accordingly, pass it the header
	handle_err := handleFile(header, args)
	check(handle_err)
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
	if args["silent"] == "false" {
		for _, v := range keys {
			fmt.Printf("%7v", v)
		}
	}
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
		usage0 = poll_cpu()
		// wait for the desired refresh time
		time.Sleep(time.Duration(refresh) * time.Second)
		// get our ending poll
		usage1 := poll_cpu()
		// we iter over the keys to ensure correct ordering
		for _, k := range keys {
			// this is how we calculate the usage
			idle = usage1[k].idle - usage0[k].idle
			total = usage1[k].total - usage0[k].total
			util = fmt.Sprintf("%.2f", 100*(float64(total)-float64(idle))/float64(total))
			// append to lines arrray
			lines = append(lines, util)

		}
		// here we convert is into a string in csv flavor
		line = strings.Join(lines[:], ",")
		// line += "\n"
		if args["silent"] == "false" {
			// if they want to see usage during runtime, give them
			// everything but the unixtime in the first position
			data = lines[1:]
			// this println gives us a newline cleanly-ish
			fmt.Println()
			for _, v := range data {
				// we pad with seven since the stat cant be greater
				// than 100.99 which has a length of 6
				fmt.Printf("%7v", v)
			}
		}
		// this tool in its current state ALWAYS creates a record
		if args["db"] != "false" {
			err := send2db(args, line)
			if err != nil {
				if retry >= 10 {
					exit_msg := "ran out of retries when connected to the DB, exiting..."
					fmt.Println(exit_msg)
					exit_err := cpuLog(exit_msg)
					check(exit_err)
					os.Exit(1)
				}
				err := cpuLog(err)
				check(err)
				time.Sleep(30 * time.Second)
				retry += 1
				info := fmt.Sprintf("encountered an error when connecting to the DB, waiting 30s and trying again (attempt %v/10)", retry)
				fmt.Println(info)
				e := cpuLog(info)
				check(e)
				continue
			}
			retry = 0
		} else {
			line += "\n"
			write_err := writeFile(args, line)
			check(write_err)
		}
	}
	// exit gracefully
	fmt.Println()
	os.Exit(0)
}
