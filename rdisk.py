import time, psutil, datetime
import regulartools as rtk

def io_poll(poll_time):
    r_actions = []
    w_actions = []
    r_bytes = []
    w_bytes = []
    count = 0
    while count <= poll_time:
        disk_start = psutil.disk_io_counters(perdisk=False)
        time.sleep(1)
        disk_end = psutil. disk_io_counters(perdisk=False)
        read_start = disk_start.read_count
        read_end =  disk_end.read_count
        write_start = disk_start.write_count
        write_end = disk_end.write_count
        r_actions_data = read_end - read_start
        w_actions_data = write_end - write_start
        r_bytes_start = disk_start.read_bytes
        r_bytes_end = disk_end.read_bytes
        w_bytes_start = disk_start.write_bytes
        w_bytes_end = disk_end.write_bytes
        r_bytes_data = r_bytes_end - r_bytes_start
        w_bytes_data = w_bytes_end - w_bytes_start
        r_actions.append(r_actions_data)
        w_actions.append(w_actions_data)
        r_bytes.append(r_bytes_data)
        w_bytes.append(w_bytes_data)
        count += 1
    return r_actions, w_actions, r_bytes, w_bytes

def crunch_data(run_limit):
    start = time.time()
    run_time = 0
    while run_time <= run_limit:
        r_count, w_count, r_bytes, w_bytes = io_poll(10)
        io_data = open('disk_io.plot', 'a', 1)
        byte_data = open('disk_bytes.plot', 'a', 1)
        human_stats = open('io_stats.txt', 'w')
        for i in range(0, len(r_count)):
            io_plot_data = str(r_count[i]) + ', ' + str(w_count[i]) + '\n'
            io_data.write(io_plot_data)
        for i in range(0, len(r_bytes)):
            byte_plot_data = str(r_bytes[i]) + ', ' + str(w_bytes[i]) + '\n'
            byte_data.write(io_plot_data)
        max_r_bytes = rtk.make_readable(max(r_bytes))
        max_w_bytes = rtk.make_readable(max(w_bytes))
        max_r_action = max(r_count)
        max_w_action = max(w_count)
        avg_r_bytes = rtk.make_readable(sum(r_bytes) / len(r_bytes))
        avg_w_bytes = rtk.make_readable(sum(w_bytes) / len(w_bytes))
        avg_r_action = format(sum(r_count) / len(r_count), '.2f')
        avg_w_action = format(sum(w_count) / len(w_count), '.2f')
        print("most bytes read ".ljust(20) + str(max_r_bytes) + '\n' +
                "average read ".ljust(20) + str(avg_r_bytes) + '\n' +
                "most bytes written ".ljust(20) + str(max_w_bytes) + '\n' +
                "average write ".ljust(20) + str(avg_w_bytes) + '\n' +
                "max read operation ".ljust(20) + str(max_r_action) + '\n' +
                "avg read operation ".ljust(20) + str(avg_r_action) + '\n' +
                "max write operation ".ljust(20) + str(max_w_action) + '\n' +
                "avg write operation ".ljust(20) + str(avg_w_action) + '\n')
        byte_data.close()
        io_data.close()
        now = time.time()
        run_time = now - start


crunch_data(12000)
