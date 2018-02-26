"""
A python tool set
"""

import random
import pickle
import copy

def sec(seconds):
    """
    A python program that converts seconds to
    human readable times
    """
    #Copies the number
    s = int(seconds)

    #Sets default values
    year = 0
    month = 0
    week = 0
    day = 0
    hour = 0
    mins = 0

    # this loop tests the input and breaks
    # it into years, months, weeks, days,
    # hours, and seconds
    while True:
        if s <= 0:
            print("sorry thats zero")
            error = 'true'
            break
        elif s >= 31556952:
            s -= 31556952
            year += 1
            continue
        elif s >= 2628000:
            s -= 2628000
            month += 1
            continue
        elif s >= 603120:
            s -= 603120
            week += 1
            continue
        elif s >= 86160:
            s -= 86160
            day += 1
            continue
        elif s >= 3600:
            s -= 3600
            hour += 1
            continue
        elif s >= 60:
            s -= 60
            mins += 1
            continue
        elif s < 60:
            break
    print(str(seconds) + " seconds  is  " + str(year)
    + "y " + str(week) + "w " + str(day)
    + "d " + str(hour) + "h " + str(mins)
    + "m " + str(s) + "s")

##
##  greeting banner
##
def greet(msg):
    """
    creates a simple banner and greeting
    """
    message = " Welcome to the " + msg.title()
    print("\n")
    print(message.center(80))
    print('_' * 80 + "\n" + '-' * 80 + "\n")


###
###  bytes to human readable
###

def make_readable(val):
    """
    a function that converts bytes via base 2 (binary)
    instead of base 10 (decimal) to human readable forms
    """
    data = int(val)
    tib = 1024 ** 4
    gib = 1024 ** 3
    mib = 1024 ** 2
    kib = 1024
    if data >= tib:
        symbol = '  TB'
        new_data = data / tib
    elif data >= gib:
        symbol = '  GB'
        new_data = data / gib
    elif data >= mib:
        symbol = '  MB'
        new_data = data / mib
    elif data >= kib:
        symbol = '  kB'
        new_data = data / kib
    else:
        symbol = '  B'
        new_data = data
    formated_data = "{0:.2f}".format(new_data)
    converted_data = str(formated_data) + symbol
    return converted_data

##
##  unit conversion
##

#-(n)by2(n)by = (n)by-#
#-by is bytes
def by2kby(val):
    x = val / 1024
    return x

def by2mby(val):
    x = val / ( 1024 ** 2 )
    return x

def by2gby(val):
    x = val / ( 1024 ** 3 )
    return x

def by2tby(val):
    x = val / ( 1024 ** 4 )
    return x

def tby2by(val):
    x = val * ( 1024 ** 4 )
    return x

def gby2by(val):
    x = val * ( 1024 ** 3 )
    return x

def mby2by(val):
    x = val * ( 1024 ** 2 )
    return x

def kby2by(val):
    x = val * 1024
    return x

def kby2mby(val):
    x = val / 1024
    return x

def kby2gby(val):
    x = val / ( 1024 ** 2 )
    return x

def kby2tby(val):
    x = val / ( 1024 ** 3 )
    return x

def tby2kby(val):
    x = val * ( 1024 ** 3 )
    return x

def gby2kby(val):
    x = val * ( 1024 ** 2 )
    return x

def mby2kby(val):
    x = val * 1024
    return x


##--substitution based encryption
##
def decode(data, key):

    ## opens files supplies by function parameter
    input_file = open(data, 'rb')
    decoder = open(key, 'rb')
    active = True
    ## while active is true
    while active:
        try:
            ## dump data into 'contents' until EOF
            contents = pickle.load(input_file)
        except EOFError:
            ## if EOF, active is false
            active = False
    ## same as above but for decoder
    active = True
    while active:
        try:
            load_key = pickle.load(decoder)
        except EOFError:
            active = False
    ## sets empty list for converting
    converted_string = []
    ## for each value in the encrypted data,
    for ch in contents[:]:
        ## if the value is in the key,
        if ch in load_key:
            ## add the decoded character to converted_string
            converted_string.append(load_key[ch])
        else:
            converted_string.append(ch)
    clear_text = ''.join(str(e) for e in converted_string)
    ## print unencrypted results
    decoded = open('decoded.txt', 'w')
    decoded.write(clear_text)
    decoded.close()
    #print(clear_text)
    input_file.close()
    decoder.close()


def encode(method, data):
    """
    a function that encodes a file with a specified method
    and creates a converted data file
    """
    user_string = []
    new_string = []
    encoder = open(method, 'rb')
    datafile = open('data.dep', 'wb')
    active = True
    ## while active is true
    while active:
        try:
            ## dump data into 'contents' until EOF
            contents = pickle.load(encoder)
        except EOFError:
            ## if EOF, active is false
            active = False

    with open(data) as infile:
        ## while true, examine each character
        while True:
            ## append 'user_string' with each character in file
            ch = infile.read(1)
            user_string.append(ch)
            ## if the character is in the conversion pool
            if ch in contents:
                ## convert character
                new_string.append(contents[ch])
            else:
                ## keep character as is, (rare)
                new_string.append(ch)
            ## repeat until EOF
            if not ch:
                eof = 'true'
                break
    pickle.dump(new_string, datafile)
    datafile.close()

def makekey(method, key):
    ## oens file to store key in
    keyfile = open(key, 'wb')
    methodfile = open(method, 'wb')
    ## defines some default values
    key = []
    working_pool = {}
    out_key = {}
    ## character map, gets copied to create character pools
    all_char = ['a', 'b', 'c', 'd', 'e',
                'f', 'g', 'h', 'i', 'j',
                'k', 'l', 'm', 'n', 'o',
                'p', 'q', 'r', 's', 't',
                'u', 'v', 'w', 'x', 'y',
                'z', 'A', 'B', 'C', 'D',
                'E', 'F', 'G', 'H', 'I',
                'J', 'K', 'L', 'M', 'N',
                'O', 'P', 'Q', 'R', 'S',
                'T', 'U', 'V', 'W', 'X',
                'Y', 'Z', '0', '1', '2',
                '3', '4', '5', '6', '7',
                '8', '9', '!', '@', '#',
                '$', '%', '^', '&', '*',
                '+', '-', '_', '=', '~',
                ' ', '.', ',', ';', '(',
                ')', '<', '>', '?', ':',
                '|', '/', '[', ']', '{',
                '}', '\'', '"', '\t',
                '\n', '\\']
    ## above list copied
    char_copy = copy.deepcopy(all_char)
    pool = copy.deepcopy(all_char)

    ## for each character in the character map,
    for char in char_copy[:]:
        ## pick a random character from pool
        temp_char = random.choice(pool)
        ## appends working_pool dictionary
        ## {'random': 'actual'} EXAMPLE
        working_pool[temp_char] = char
        ## appends my_key dictionary, creates key
        ## {'actual': 'random'} EXAMPLE
        out_key[char] = temp_char
        ## remove the random character assigned
        ## from the pool so it cant be picked twice
        pool.remove(temp_char)
        pickle.dump(out_key, keyfile)
        pickle.dump(working_pool, methodfile)
