import psutil as ps
import rextools as rtk
import time
import os

def net_poll(poll_time):
    sys_wide_start = ps.net_io_counters(pernic=False)
    per_nic_start = ps.net_io_counters(pernic=True)
    time.sleep(poll_time)
    sys_wide_end = ps.net_io_counters(pernic=False)
    per_nic_end = ps.net_io_counters(pernic=True)
    return sys_wide_start, per_nic_start, sys_wide_end, per_nic_end


def crunch_data():
    sys_start, by_nic_start, sys_end, by_nic_end  = net_poll(30)
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
    total_r_plot = open('total_recv.plot', 'a', 1)
    total_s_plot = open('total_sent.plot', 'a', 1)
    err_in_plot = open('error_in.plot', 'a', 1)
    err_out_plot = open('error_out.plot', 'a', 1)
    human_stats = open('netstats.txt', 'w')
    human_stats.write("Total Sent".ljust(18) + str(total_sent) + '\n'
                      + "Total Received".ljust(18) + str(total_recv) + '\n'
                      + "Dropped in/out".ljust(18) + str(error_in) + '/'
                      + str(error_out) + '\n')
    human_stats.close()
    total_r_plot.write(str(recv_data) + '\n')
    total_r_plot.close()
    total_s_plot.write(str(sent_data) + '\n')
    total_s_plot.close()
    err_in_plot.write(str(error_in) + '\n')
    err_in_plot.close()
    err_out_plot.write(str(error_out) + '\n')
    err_out_plot.close()


while True:
    crunch_data()
    f = open('netstats.txt')
    text =f.read()
    print(text)
    f.close()
