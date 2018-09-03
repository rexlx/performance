import psutil as ps
import regulartools as rtk
import time, os, sys

def net_poll(poll_time):
    sys_wide_start = ps.net_io_counters(pernic=False)
    per_nic_start = ps.net_io_counters(pernic=True)
    time.sleep(poll_time)
    sys_wide_end = ps.net_io_counters(pernic=False)
    per_nic_end = ps.net_io_counters(pernic=True)
    return sys_wide_start, per_nic_start, sys_wide_end, per_nic_end


def crunch_data():
    sys_start, by_nic_start, sys_end, by_nic_end  = net_poll(5)
    recv_start = sys_start.bytes_recv
    recv_end =  sys_end.bytes_recv
    sent_start = sys_start.bytes_sent
    sent_end = sys_end.bytes_sent
    error_in = sys_end.dropin
    error_out = sys_end.dropout
    recv_data = recv_end - recv_start
    sent_data = sent_end - sent_start
    total_recv = rtk.make_readable(recv_data)
    total_sent = rtk.make_readable(sent_data)
    now = str(time.time())
    with open('net.plot', 'a') as f:
        f.write(now + ', ' + str(recv_data)  + ', ' + str(sent_data)  + ', ' +
                str(error_in)  + ', ' + str(error_out) + '\n')

    msg = "sent/recv: " + str(total_sent) + '/' + str(total_recv) + \
    '  errors i/o: ' + str(error_in) + ' / ' + str(error_out)
    sys.stdout.write('%s\r' % msg)
    sys.stdout.flush()

def main():
    runtime = 5760
    c = 0
    while c < runtime:
        crunch_data()
        c += 1

if __name__ == '__main__':
    main()
