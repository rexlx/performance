import time, psutil, datetime
import regulartools as rtk

def io_poll(poll_time):
    io_start = psutil.disk_io_counters(perdisk=False)
    time.sleep(poll_time)
    io_end = psutil.disk_io_counters(perdisk=False)
    return io_start, io_end

def crunch_data():
    disk_start, disk_end = io_poll(5)
    read_start = disk_start.read_count
    read_end =  disk_end.read_count
    write_start = disk_start.write_count
    write_end = disk_end.write_count
    read_data = read_end - read_start
    write_data = write_end - write_start
    data = str(read_data) + ', ' + str(write_data) + '\n'
    ## will add by byte counters soon
    #total_read = rtk.make_readable(read_data)
    #total_wrote = rtk.make_readable(write_data)
    io_data = open('disks.plot', 'a', 1)
    human_stats = open('io_stats.txt', 'w')
    human_stats.write("Read ".ljust(18) + str(read_data) + '\n'
                      + "Write".ljust(18) + str(write_data) + '\n')
    human_stats.close()
    io_data.write(data)
    io_data.close()

run_time = 0
while run_time < 640:
    crunch_data()
    f = open('io_stats.txt')
    text = f.read()
    print(text)
    f.close()
    run_time += 1
