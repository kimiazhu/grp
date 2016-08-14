#!/bin/bash

#=================================================================
# GRP简单部署脚本。
#
# Author: ZHU HAIHUA
# Version: 1.0 创建脚本。(2016-08-14)
# Since: 2016-08-14
#=================================================================

FILENAME="./grp"
PROC_NAME="grp"

function restart() {
    echo "ready to restart"
    pgrep ${PROC_NAME} | xargs kill -9
    mv -f ${FILENAME} ${PROC_NAME}
    chmod a+x ${FILENAME}
    nohup ${FILENAME} > /dev/null 2>stderr.log &
    sleep 1
    echo "finish restart"
    pgrep ${PROC_NAME}
}

function reload() {
    echo "ready to reload"
    pgrep ${PROC_NAME} | xargs kill -HUP
    sleep 1
    echo "finish reload"
    pgrep ${PROC_NAME}
}

function kill() {
    echo "ready to kill grp"
    pgrep ${PROC_NAME} | xargs kill -9
    echo "finish"
}

function reload2() {
    echo "ready to reload"
    pids=`pgrep ${PROC_NAME}`
    for pid in $pids
    do
        #ssh -p$PORT $SERVER "kill -HUP $pid"
        echo "send hup signal to remote process: $pid"
    done
}

function usage() {
IFS=
usage_msg="
GRP deploy script.
Userage: grp.sh [-h] [-r|-s|-k]
    -h, --help    Show usage
    -r, --restart [grp_file] restart GRP, if <grp_file> exists, it will be rename to grp and then restart
    -s, --reload reload GRP
    -k, --kill stop current GRP service

Examples:
    # stop grp, rename grp_2 to grp then restart
    $ ./grp.sh -r grp_2
    # reload grp
    $ ./grp.sh -s
"
echo $usage_msg
IFS=,
}

#=======================================================================
########################### MAIN PROCCESS ##############################
#=======================================================================
ARGS=`getopt -o hr::sk --long help,restart::,reload,kill -n 'grp.sh' -- "$@"`
if [ $# -eq 0 ]; then
    echo "Terminating... see -h or --help for help"
    exit 1
fi
eval set -- "${ARGS}"

while [[ true ]]; do
    case "$1" in
        -h|--help)
            shift
            SHOW_HELP=yes
            ;;
        -r | --restart)
            RESTART=yes
            case $2 in
                "")
                    echo "restart grp process"
                    shift 2
                    ;;
                *)
                    FILENAME=$2
                    echo "restart using ${FILENAME}"
                    shift 2
                    ;;
            esac
            ;;
        -s|--reload)
            shift
            RELOAD=yes
            ;;
        -k | --kill)
            shift
            KILL=yes
            ;;
        --)
            shift
            #break
            ;;
        -*)
            echo "Wrong option~~ Please read the flowing usage doc."
            echo ""
            usage
            exit 2
            ;;
        *)
            break
            ;;
    esac
done

if [ "$SHOW_HELP" = "yes" ]; then
    usage
    exit 0
fi

if [ -n "${RELOAD}" ]; then
    reload
    exit 0
fi

if [ -n "${RESTART}" ]; then
    restart
    exit 0
fi

if [ -n "${KILL}" ]; then
    kill
    exit 0
fi
