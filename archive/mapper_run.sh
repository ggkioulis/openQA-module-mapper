#! /bin/bash

now=`date +%Y%m%d%H%M%S`;
/root/github/openQA-module-mapper/archive/scrapper &> /root/github/openQA-module-mapper/archive/logs/$now.log 
