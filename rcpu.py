from __future__ import print_function
from __future__ import division
import time, psutil, sys, platform, os, argparse

# runs for eight hours, one poll per second, screen updates every 5 seconds
#    refresh = 5
#    run_time = 5760
# 60 seconds / refresh (5) = 12; it takes one minute to perform 12 refreshes
refresh = 5
run_time = 5760

def get_args():
    # create parser
        msg = "This script records network statistics"
        parser = argparse.ArgumentParser(description=msg)
        # add expected arguments
        parser.add_argument('-s', dest='silent', required=False,
                            action="store_true", help="display statistics")
        parser.add_argument('-n', dest='noheader', required=False,
                            action="store_true", help="dont write header")
        parser.add_argument('-R', dest='refresh', required=False)
        parser.add_argument('-r', dest='runtime', required=False)
        args = parser.parse_args()
        if args.silent:
            silent = True
        else:
            silent = False
        if args.noheader:
            noheader = True
        else:
            noheader = False
        if args.refresh:
            refresh = float(args.refresh)
        else:
            # default refresh i s 5 seconds
            refresh = 5
        if args.runtime:
            runtime = float(args.runtime)
        else:
            # default runtime is eight hours
            runtime = 28800
        return silent, noheader, refresh, runtime


def get_dist():
    """
    this function gets the os and determines if the machine
    is real or virtual
    """
    cpuinfo = '/proc/cpuinfo'
    poll_type = platform.system()
    if poll_type == 'Windows':
        return poll_type
    # this block tests /proc/cpuinfo for the hypervisor flag
    elif poll_type == 'Linux':
        return poll_type
    elif os.path.isfile(cpuinfo):
        try:
            with open(cpuinfo) as f:
                for line in f:
                    if 'hypervisor' in line:
                        poll_type = 'hypervisor'
                        return poll_type
        except Exception as e:
            print('couldnt open /proc/cpuinfo!')
            sys.exit()


def my_poll(refresh, poll_type):
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
    while count < refresh:
        if poll_type == 'Windows' or poll_type == 'hypervisor':
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
        # dump data into a csv(,)
        #epoch, cpu load, frequency(linux only)
        with open(outfile, 'a') as f:
            for i in range(0, len(loads)):
                f.write(str(times[i]) + ',' + str(loads[i]) +
                        ',' + str(speeds[i]) + '\n')
        return cpu_avg, min_speed, max_speed
    else:
        # epoch, cpu load
        with open('cpuutil.plot', 'a') as f:
            for i in range(0, len(loads)):
                f.write(str(times[i]) + ',' + str(loads[i]) + '\n')
        return cpu_avg

def main(silent, noheader, refresh, runtime):
    poll_type = get_dist()
    uptime = 0
    start = time.time()
    try:
        # we dont want to append an old file, remove it
        os.remove('cpuutil.plot')
        # if they want a header
        if not noheader:
            if poll_type == 'Linux':
                with open('cpuutil.plot', 'w') as f:
                    f.write('utime,load,speed\n')
            else:
                with open('cpuutil.plot', 'w') as f:
                    f.write('utime,load\n')
    except OSError:
        pass
    while uptime <= runtime:
        if poll_type != 'Windows' and poll_type != 'hypervisor':
            load, min_freq, max_freq = my_poll(refresh, poll_type)
            data = "average usage: " + str(round(load, 2)) + \
                   " cpu speed min/max: " + str(round(min_freq, 2)) + '/' + \
                   str(round(max_freq, 2))
            if not silent:
                sys.stdout.write('%s\r' % data)
                sys.stdout.flush()
        else:
            load = my_poll(refresh, poll_type)
            data = format(load, '.2f')
            msg = "percent used:  " + str(data) + ' '
            if not silent:
                sys.stdout.write('%s\r' % msg)
                sys.stdout.flush()
        now = time.time()
        uptime = now - start

if __name__ == '__main__':
    silent, noheader, refresh, runtime = get_args()
    main(silent, noheader, refresh, runtime)
    print('\n')
