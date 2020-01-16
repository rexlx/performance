# performance
A collection of performance gathering tools

You'll notice ive split this into two categories, python and golang.
golang currently has limited support (only rcpu is ready), rmem will
come next (mostly done). The rest of the readme only covers python (golang will get its own)

These tools rely on python and psutil. If you plan to plot the data,
matplotlib should work fine. All data is currently written to a
respective csv and can be plotted easily.

# Golang

**1. rcpu**
<br>
collects cpu stats and can either write them to stdout while recording them to a csv or sending them to a mongodb instance.


run without displaying to stdout refreshing every 5s forever
```bash
rcpu -r inf -R 5 -d "mongodb://192.168.1.42:9091" -s
```

display stats to screen while writing to csv, run for eight hours
```bash
rcpu -r 28800 -R 5
```

**2. rmem**
<br>
coming soon

# Python
**1. rcpu.py**

gathers the cpu load and, if you're running linux, frequency. frequency is collected as a single value, and is the average of all core speeds. currently configured to write plot data to ./cpuutil.plot as csv(,). fields are:<br/>
*unixtime,cpu load,frequency(linux only)*


**2. rnet.py**

gathers network stats and writes them to a file ./net.plot with the following structure:

*utime,recv,sent,err_in,err_out*


**3. rdisk.py**

gathers disk stats from each drive on the system. saved as DISKNAME.plot, fields are:<br/>
*unixtime, reads, writes, bytes read, bytes written, read wait, write wait*
**NOTE** on newer kernels that have added new fields to iostat, this will break. to fix:
```
pip uninstall psutil
pip install git+https://github.com/giampaolo/psutil.git
```

**4. rmem.py *linux only***<br/>

gathers memory and swap stats<br/>
*utime,total,used,free,buff,cache,slab,swap*<br/>

**5. total_rw.py**

shows amount written / read since current boot. does not plot.



**6. get_cpu_speed**

will collect the current cpu speed (in mhz) for each core (including hyperthreaded ones) once a second for the length of the poll(SECONDS) function. creates a csv named cpu_speeds.csv with the following headers:<br/>
*unixtime,cpu0,cpu2,cpu3,etc...,avg*

**NOTE** only works on non-virtualized unix based OS'
