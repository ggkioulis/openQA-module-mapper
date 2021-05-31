#! /bin/bash

go build ..;
now=`date +%Y%m%d%H%M%S`;
/root/github/openQA-module-mapper/archive/openQA-module-mapper &> /root/github/openQA-module-mapper/archive/logs/$now.log 
