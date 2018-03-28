import psutil as ps
import time, datetime

def one_min():
    one_min_load = []
    count = 0
    while count <= 60:
        meminfo = ps.virtual_memory()
        total_mem = meminfo.total
        usage = meminfo.percent
        time.sleep(1)
        one_min_load.append(usage)
        count += 1
    with open('memutil.plot', 'a', 1) as loadfile:
        for load in one_min_load:
            loadfile.write(str(load) + '\n')
    max_int = max(one_min_load)
    return max_int

run_time = 0
while run_time < 640:
    load = one_min()
    print("max interval for poll period: " + str(load) + '\n')
    run_time += 1
