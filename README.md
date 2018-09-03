# performance
A collection of performance gathering tools

These tools rely on python and psutil. If you plan to plot the data,
matplotlib should work fine. All data is currently written to a
respective csv and can be plotted easily.



1. rcpu.py:

gathers the cpu load and, if you're running linux, frequency.
currently configured to write plot data to ./cpuutil.plot as csv(,).
epoch, cpu load, frequency(linux only)


2. rnet.py:

gathers network stats over a specified poll time untill the process is stopped
creates the following files:

error_out.plot, error_in.plot, total_sent.plot, total_recv.plot

3. rdisk.py

gathers i/o stats and stores the values 'R'    'W' separated by tab in a file
disks.plot
