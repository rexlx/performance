from __future__ import print_function
from __future__ import division
import time, psutil, time, sys, platform

# runs for eight hours, one poll per second, screen updates every 5 seconds
#    refresh = 5
#    run_time = 5760
# 60 seconds / refresh (5) = 12; it takes one minute to perform 12 refreshes
refresh = 5
run_time = 5760

def get_dist():
    """
    this function gets the OS
    """
    os = platform.system()
    if os == 'Windows':
        return os
    else:
        dist = platform.linux_distribution()
        os = ' '.join(e for e in dist)
        return os

def my_poll(rate, os):
    """
    this function gathers the cpu usage and frequency (freq doesnt work on
    windows) once a second per refresh period(rate) and updates the screen
    once every refresh period. it then writes all per second poll / freq data
    to a csv file.
    """
    outfile = 'cpuutil.plot'
    times = []
    loads = []
    speeds = []
    count = 0
    while count < rate:
        if os == 'Windows':
            # mark the time
            now = str(time.time())
            times.append(now)
            # get the cpu usage
            x = psutil.cpu_percent(interval=None, percpu=False)
            # add it to the load list
            loads.append(x)
        else:
            # mark the time
            now = str(time.time())
            times.append(now)
            # get the cpu usage
            x = psutil.cpu_percent(interval=None, percpu=False)
            # add it to the load list
            loads.append(x)
            # get the cpu speed
            freq = psutil.cpu_freq()
            # i only care about the current speed
            cpu_speed = freq.current
            # append the speeds list
            speeds.append(cpu_speed)
        # wait a second and increment. repeat.
        time.sleep(1)
        count += 1
    # determine the average usage for period
    cpu_avg = sum(loads) / len(loads)
    # determine low / high speed for poll period
    if len(speeds) > 0:
        min_speed, max_speed = min(speeds), max(speeds)
        # dump the data into a csv file where the first column is cpu util
        # and the second is the current cpu speed
        with open(outfile, 'a') as f:
            for i in range(0, len(loads)):
                f.write(str(times[i]) + ', ' + str(loads[i]) +
                        ', ' + str(speeds[i]) + '\n')
        return cpu_avg, min_speed, max_speed
    else:
        #only write the per second values into a file
        with open('cpuutil.plot', 'a') as f:
            for i in range(0, len(loads)):
                f.write(str(times[i]) + ', ' + str(loads[i]) + '\n')
        return cpu_avg

def main(run_time):
    c = 0
    while c < run_time:
        os = get_dist()
        if os != 'Windows':
            load, min_freq, max_freq = my_poll(refresh, os)
            data = "average usage: " + str(round(load, 2)) + \
                   " cpu speed min/max: " + str(round(min_freq, 2)) + '/' + \
                   str(round(max_freq, 2))
            sys.stdout.write('%s\r' % data)
            sys.stdout.flush()
        else:
            load = my_poll(refresh, os)
            data = format(load, '.2f')
            msg = "percent used:  " + str(data) + ' '
            sys.stdout.write('%s\r' % msg)
            sys.stdout.flush()
        c += 1

if __name__ == '__main__':
    main(run_time)
