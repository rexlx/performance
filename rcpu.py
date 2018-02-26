import time, psutil, datetime

time_to_run = 640
print('Collecting CPU stats for ' + str(time_to_run))
def one_min():
    one_min_load = []
    count = 0
    while count <= 60:
        x = psutil.cpu_percent(interval=None, percpu=False)
        one_min_load.append(x)
        time.sleep(1)
        count += 1
    one_min_avg = sum(one_min_load) / len(one_min_load)
    with open('uptime.plot', 'a') as loadfile:
        for load in one_min_load:
            loadfile.write(str(load) + '\n')
    return one_min_avg

run_time = 0
outfile = open('uptime.txt', 'a', 1)
timer = open('runtime.txt', 'a', 1)
now = datetime.datetime.now()
hour, minute = str(now.hour), str(now.minute)
the_time = hour + ':' + minute
timer.write(the_time + '\n')
while run_time < 640:
    load = one_min()
    outfile.write("%.2f" % load + '\n')
    if run_time % 10 == 0:
        update = time_to_run - run_time
        print(str(update) + 'minutes remaining\n')
    if run_time % 60 == 0:
        now = datetime.datetime.now()
        hour, minute = str(now.hour), str(now.minute)
        the_time = hour + ':' + minute
        timer.write(the_time + '\n')
    run_time += 1
outfile.close()
