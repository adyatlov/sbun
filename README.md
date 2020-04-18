# sbun
Tool for analyzing [DC/OS service diagnostics bundle](https://support.d2iq.com/s/article/create-service-diag-bundle)

Usage
-----

```
$ cd <service diagnostics bundle directory>
$ sbun
```

Features
--------

Writes service task list to the standard output in the CVS format. The order of columns is:

1. task name
1. starting timestamp
1. running timestamp
1. killing timestamp
1. task ID
1. path to the task directory
