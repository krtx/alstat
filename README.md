# alstat

[![CircleCI](https://circleci.com/gh/krtx/alstat.svg?style=svg)](https://circleci.com/gh/krtx/alstat)

Provide access summary of your ltsv log file.

## Description

`alstat` parse ltsv files and show the summary. 

## Example

```
$ cat test.log
method:GET	status:404	path: /status?id=79
method:GET	status:404	path: /profile?id=53
method:GET	status:200	path: /profile?id=90
method:POST	status:404	path: /status?id=86
method:POST	status:200	path: /profile?id=2
...
```

Specify a label to show access counts.

```
$ alstat -c 0 -l 'method' ./test.log
method  access
--------------
GET         40
POST        60
```

Labels can be combined.

```
$ alstat -c 0 -l 'method' -l 'status' ./test.log
method  status  access
----------------------
GET     200         18
GET     404         22
POST    200         24
POST    404         36
```

As an advanced usage, regexps can be specified to extract a part of
the value.

```
$ alstat -c 0 -l 'path:(/profile|/status)' ./test.log
path      access
----------------
/profile      59
/status       41
```

The first label can be used as a primary label: `-s` separates them
and `-r` prints rates for each lines within the primary label.

```
$ alstat -s -r -c 0 -l 'path:(/profile|/status)' -l status ./test.log
path      status  access   (rate)
---------------------------------
/profile  200         24   40.68%
/profile  404         35   59.32%
---------------------------------
/status   200         18   43.90%
/status   404         23   56.10%
```

`-c` specify the display interval in seconds and `-c 0` means to print
just once.  The default value is 1.

```
$ alstat -s -r -l 'path:(/profile|/status)' -l status ./test.log
... access summary is displayed on every second ...
```
