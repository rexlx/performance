import regulartools as rtk

disks = {}

try:
    with open('/proc/diskstats', 'r') as f:
        for line in f:
            name, amount = line.split()[2], line.split()[9]
            disks[name] = amount
            can_move_on = True
except IOError:
    print("\ncould not find '/proc/diskstats'!")
    can_move_on = False

if can_move_on:
    for k, v in disks.items():
        bytes = int(v) * 512
        total_written = rtk.make_readable(bytes)
        print(k.ljust(16) + total_written)
