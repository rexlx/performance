import os, time
import multiprocessing as mp

"""
description:

this script opens /proc/cpuinfo and extracts the current mhz as an
integer. then, after adding a header line, writes the data to a csv.
currently the runtimme is hardcoded into the poll function and can only
be mofified there. i plan to add some args for runtime in seconds,
delimiter type (currently ',') and who knows what else :), the cpu
speed can only be obtained on bare metal machines (not virtualized),
and even then, some machines have bios features that may control the
speed and therefore obscure that information to the OS.
coughlenovocough.

example:

python get_cpu_speed.py

creates a file 'cpu_speeds.csv' (need and arg for that too)

utime,cpu0,cpu1,cpu2,cpu3
1544279626.23,2749,3005,2429,2889
"""


def cpu_speeds():
    """
    gets the current cpu speed for each processor via /proc/cpuinfo
    returns unixtime and cpu speeds as python list
    """
    # where the cpu info is
    with open('/proc/cpuinfo') as f:
        # create an empty list for later
        speeds = []
        for line in f:
            # get the unix time via time module
            utime = time.time()
            # we only care about the cpu frequency line
            if 'MHz' in line:
                # split the line up into indexes
                x = line.split()
                # convert the indexed string to a float and then to int
                # to shed the decimals
                speed = int(float(x[3]))
                # add it to the list
                speeds.append(str(speed))
    return utime, speeds

def make_csv(utime, speeds):
    """
    this takes two parameters, utime and speeds, writes data to a csv
    named cpu_speeds.csv. the rows are as follows:
    unixtime,cpu0,cpu1,cpu2,cpu3,etc.....
    """
    with open('cpu_speeds.csv', 'a') as f:
        # write the unixtime (no newline))
        f.write(str(utime) + ',')
        # write all cpu data except the last one + ,
        for i in speeds[0:-1]:
            f.write(str(i) + ',')
        # write the final cpu data
        f.write(speeds[-1])
        # finally, add the newline
        f.write('\n')

def poll(runtime):
    """
    this function adds the header bar to the csv, gets the cpu speeds, and
    then writes it to a csv. requires runtime in seconds as parameter
    """
    # get the total cores (including hyperthreaded ones)
    total_cores = mp.cpu_count()
    # mark the start time
    then = time.time()
    # initialize length for the while loop below
    length = 0
    # list comprehension that concatenates 'cpu' + core number + ','
    # calculated from the range function 0-total_cores
    cpus = ['cpu' + str(e) + ',' for e in range(0, total_cores)]
    # we dont want the last comma
    cpus[-1] = cpus[-1].rstrip(',')
    # if the csv exists, remove it
    try:
        os.remove('cpu_speeds.csv')
    except OSError:
        pass
    with open('cpu_speeds.csv', 'a') as fname:
        # write the csv header
        fname.write('utime,')
        for i in cpus:
            fname.write(i)
        fname.write('\n')
    # main loop
    while length <= runtime:
        # unpack the return of cpu_speeds
        utime, speeds = cpu_speeds()
        # feed that into make_csv
        make_csv(utime, speeds)
        # wait
        time.sleep(1)
        # determine runtime
        now = time.time()
        length = now - then
        
if __name__ == "__main__":
    poll(28800)
