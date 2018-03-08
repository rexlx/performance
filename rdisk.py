import time, psutil, datetime

time_to_run = 640
print('Collecting disk stats for ' + str(time_to_run))
def one_min():
    one_min_read = []
    one_min_write = []
    count = 0
    while count <= 2:
        x = psutil.disk_io_counters(perdisk=False)
        one_min_read.append(x.read_count)
        one_min_write.append(x.write_count)
        time.sleep(1)
        count += 1
    read_avg = sum(one_min_read) / len(one_min_read)
    write_avg = sum(one_min_write) / len(one_min_write)
    with open('disks.plot', 'a') as loadfile:
        index = 0
        while index < len(one_min_read):
            loadfile.write(str(one_min_read[index]) + ', '
                           + str(one_min_write[index]) + '\n')
            index += 1
    return read_avg, write_avg

run_time = 0
outfile = open('disks.txt', 'a', 1)
timer = open('disks_runtime.txt', 'a', 1)
now = datetime.datetime.now()
hour, minute = str(now.hour), str(now.minute)
the_time = hour + ':' + minute
timer.write(the_time + '\n')
while run_time < 640:
    read_avg, write_avg = one_min()
    read = format(read_avg, '.2f')
    write = format(write_avg, '.2f')
    data = read + ', ' + write + '\n'
    outfile.write(data)
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
