import time, psutil, datetime

def one_min():
    one_min_load = []
    count = 0
    while count <= 60:
        x = psutil.cpu_percent(interval=None, percpu=False)
        one_min_load.append(x)
        time.sleep(1)
        count += 1
    one_min_avg = sum(one_min_load) / len(one_min_load)
    with open('cpuutil.plot', 'a') as loadfile:
        for load in one_min_load:
            loadfile.write(str(load) + '\n')
    return one_min_avg

run_time = 0
while run_time < 3600:
    load = one_min()
    print("average for poll period: " + str(load) + '\n')
    run_time += 1
