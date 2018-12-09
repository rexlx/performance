# performance
A collection of performance gathering tools

These tools rely on python and psutil. If you plan to plot the data,
matplotlib should work fine. All data is currently written to a
respective csv and can be plotted easily.



**1. rcpu.py**

gathers the cpu load and, if you're running linux, frequency. frequency is collected as a single value, and is the average of all core speeds. currently configured to write plot data to ./cpuutil.plot as csv(,). fields are:
*unixtime, cpu load, frequency(linux only)*


**2. rnet.py**

gathers network stats and writes them to a file ./net.plot with the following structure:

*recv,sent,err_in,err_out<br/>
1544368488.51,308,0,0,0*


**3. rdisk.py**

gathers disk stats from each drive on the system. saved as DISKNAME.plot, fields are:
*unixtime, reads, writes, bytes read, bytes written, read wait, write wait*
**NOTE** on newer kernels that have added new fields to iostat, this will break. to fix:
```
pip uninstall psutil
pip install git+https://github.com/giampaolo/psutil.git
```
**4. total_rw.py**

shows amount written / read since current boot. does not plot.



**5. get_cpu_speed**

will collect the current cpu speed (in mhz) for each core (including hyperthreaded ones) once a second for the length of the poll(SECONDS) function. creates a csv named cpu_speeds.csv with the following headers:
*unixtime,cpu0,cpu2,cpu3,etc...*

**NOTE** only works on non-virtualized unix based OS'
