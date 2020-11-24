#!/bin/bash -v 
ap issues
ap hw
ap sw
for i in $(seq 1 3); do ssh node011${i} 'hostname; cat /proc/cpuinfo |grep clock|uniq -c'; done
for i in $(seq 1 3); do ssh node011${i} 'hostname; /bin/systemctl status tuned | grep Active'; done