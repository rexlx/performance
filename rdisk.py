import psutil
import time
import argparse

"""
description:

this script records some basic disk stats for each disk on the system
it will then create a DISK_NAME.plot file for each disk (and partition)
in the form of a csv.

it may make sense in the future to write one plot file, and dynamically
label the headers with a disk name, but for people who dont like
headers, this would pose a problem.

example:

$ python rdisk.py -h
usage: rdisk.py [-h] [-n] [-r RUNTIME]

This script records disk statistics

optional arguments:
  -h, --help  show this help message and exit
  -n          dont write header
  -r RUNTIME

sample plot file contents:

utime,read,write,rbytes,wbytes,rwait,wwait
1544479434.33,2460,0,644808704,0,846,0
"""


def get_args():
    """
    get the cli args via the argparse module
    """
    msg = "This script records disk statistics"
    # create an instance of parser with our unique msg
    parser = argparse.ArgumentParser(description=msg)
    # add expected arguments
    parser.add_argument('-n', dest='noheader', required=False,
                        action="store_true", help="dont write header")
    parser.add_argument('-r', dest='runtime', required=False)
    args = parser.parse_args()
    if args.noheader:
        noheader = True
    else:
        noheader = False
    if args.runtime:
        runtime = float(args.runtime)
    else:
        # default runtime is eight hours
        runtime = 28800
    return noheader, runtime


def write_headers():
    """
    write the column names to our csv
    """
    disks = psutil.disk_io_counters(perdisk=True)
    disk_names = list(disks.keys())
    for disk in disk_names:
        filename = disk + '.plot'
        with open(filename, 'w') as f:
            f.write('utime,read,write,rbytes,wbytes,rwait,wwait\n')


def io_poll(runtime):
    """
    this function gets the disk stats for all disks on the system
    """
    start = time.time()
    uptime = 0
    while uptime <= runtime:
        now = time.time()
        # start values
        disk_start = psutil.disk_io_counters(perdisk=True)
        time.sleep(1)
        # end values
        disk_end = psutil.disk_io_counters(perdisk=True)
        # disk names are the keys of these dictionaries
        disk_names = list(disk_start.keys())
        for name in disk_names:
            # get the start / end values for that disk
            dstart = disk_start[name]
            dend = disk_end[name]
            read_start = dstart.read_count
            read_end = dend.read_count
            write_start = dstart.write_count
            write_end = dend.write_count
            r_actions = read_end - read_start
            w_actions = write_end - write_start
            r_bytes_start = dstart.read_bytes
            r_wait_start = dstart.read_time
            r_bytes_end = dend.read_bytes
            r_wait_end = dend.read_time
            w_bytes_start = dstart.write_bytes
            w_wait_start = dstart.write_time
            w_bytes_end = dend.write_bytes
            w_wait_end = dend.write_time
            r_bytes = r_bytes_end - r_bytes_start
            r_wait = r_wait_end - r_wait_start
            w_wait = w_wait_end - w_wait_start
            w_bytes = w_bytes_end - w_bytes_start
            # file name is disk name plus .plot
            filename = name + '.plot'
            # utime,reads,writes,bytes read,bytes written,read wait,write wait
            with open(filename, 'a', 1) as f:
                f.write(str(now) + ',' + str(r_actions) + "," + str(w_actions)
                        + "," + str(r_bytes) + "," + str(w_bytes) + ","
                        + str(r_wait) + "," + str(w_wait) + "\n")
            uptime = now - start


if __name__ == '__main__':
    noheader, runtime = get_args()
    if not noheader:
        write_headers()
    io_poll(runtime)
