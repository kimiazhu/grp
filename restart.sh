# Author: ZHU HAIHUA
# Date: 8/10/16
#!/usr/bin/env bash

pgrep grp | xargs kill -9

nohup ./grp > /dev/null 2>stderr.log &

