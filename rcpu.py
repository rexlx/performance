import time, psutil, datetime

def one_min_poll():
    """
    this function gathers the cpu usage in percent value
    """
    one_min_load = []
    count = 0
    # while the runtime is less than a minute
    while count <= 60:
        # get the cpu usage
        x = psutil.cpu_percent(interval=None, percpu=False)
        # add it to the load list
        one_min_load.append(x)
        # wait a second and increment. repeat.
        time.sleep(1)
        count += 1
    # determine the average usage for that minute
    one_min_avg = sum(one_min_load) / len(one_min_load)
    # dump the data into a file (60 one second)
    with open('cpuutil.plot', 'a') as loadfile:
        for load in one_min_load:
            loadfile.write(str(load) + '\n')
    # return the average for that minute
    return one_min_avg

# this is a simple example of how to collect data
# for a certain period of time (480 minutes = 8 hours)
run_time = 0
while run_time < 480:
    load = one_min_poll()
    print("average for poll period: " + str(load) + '\n')
    run_time += 1
