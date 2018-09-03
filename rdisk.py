import time, psutil, time
import regulartools as rtk

def io_poll(poll_time):
    """
    this function gets the disk stats for all disks on the system
    """
    count = 0
    while count <= poll_time:
        now = str(time.time())
        # start values
        disk_start = psutil.disk_io_counters(perdisk=True)
        time.sleep(1)
        # end values
        disk_end = psutil.disk_io_counters(perdisk=True)
        # disk names are the keys of these dictionaries
        disk_names = list(disk_start.keys())
        for name in disk_names:
            # get the start / end values for that disk
            dstart = disk_start[name]
            dend = disk_end[name]
            read_start = dstart.read_count
            read_end =  dend.read_count
            write_start = dstart.write_count
            write_end = dend.write_count
            r_actions = read_end - read_start
            w_actions = write_end - write_start
            r_bytes_start = dstart.read_bytes
            r_wait_start = dstart.read_time
            r_bytes_end = dend.read_bytes
            r_wait_end = dend.read_time
            w_bytes_start = dstart.write_bytes
            w_wait_start = dstart.write_time
            w_bytes_end = dend.write_bytes
            w_wait_end = dend.write_time
            r_bytes = r_bytes_end - r_bytes_start
            r_wait = r_wait_end - r_wait_start
            w_wait = w_wait_end - w_wait_start
            w_bytes = w_bytes_end - w_bytes_start
            # file name is disk name plus .plot
            filename =  name + '.plot'
            #epoch, reads, writes, bytes read, bytes written, read wait, write w
            with open(filename, 'a', 1) as f:
                f.write(now + ', ' + str(r_actions) + "," + str(w_actions) +
                "," + str(r_bytes) + "," + str(w_bytes) + "," + str(r_wait) +
                "," + str(w_wait) + "\n")
        count += 1

if __name__ == '__main__':
    io_poll(3600)
