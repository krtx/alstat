# alstat

[![CircleCI](https://circleci.com/gh/krtx/alstat.svg?style=svg)](https://circleci.com/gh/krtx/alstat)
[![Build Status](http://localhost:8000/api/badges/krtx/alstat/status.svg)](http://localhost:8000/krtx/alstat)

Provide access summary of your ltsv log file.

## Description

`alstat` parses the tail lines of given ltsv files and prints the summary.

## Example

```
$ cat test.log
method:POST	status:200	path:/profile?id=49	reqtime_microsec:7583
method:POST	status:404	path:/profile?id=4	reqtime_microsec:8931
method:GET	status:404	path:/status?id=40	reqtime_microsec:1735
method:POST	status:404	path:/profile?id=10	reqtime_microsec:9546
method:GET	status:200	path:/status?id=77	reqtime_microsec:9515
...
```

Use `-l` to specify a label to show access counts.

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

Regexps can be used to extract a part of the values.

```
$ alstat -c 0 -l 'path:(/profile|/status)' ./test.log
path      access
----------------
/profile      59
/status       41
```

The first label can be used as a primary label: `-sep` separates them
and `-rate` prints rates for each lines within the primary label.

```
$ alstat -sep -rate -c 0 -l 'path:(/profile|/status)' -l status ./test.log
path      status  access   (rate)
---------------------------------
/profile  200         24   40.68%
/profile  404         35   59.32%
---------------------------------
/status   200         18   43.90%
/status   404         23   56.10%
```

`-sum` sums up the specified fields.

```
$ alstat -sep -rate -c 0 -l method -l status -sum reqtime_microsec test.log
method  status  access   (rate)  sum(reqtime_microsec)
------------------------------------------------------
GET     200         31   58.49%                 155181
GET     404         22   41.51%                 125124
------------------------------------------------------
POST    200         21   44.68%                  95786
POST    404         26   55.32%                 134721
```

`-c` specifies the display interval in seconds and `-c 0` means to print
just once.  The default value is 1.

```
$ alstat -s -r -l 'path:(/profile|/status)' -l status ./test.log
... access summary is displayed on every second ...
```
