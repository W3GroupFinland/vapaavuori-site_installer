#! /bin/bash

# /etc/init.d/sitestatusd
# Init script for sitestatus
#
# chkconfig: 2345 95 05
# description: Sitestatus program daemon

### BEGIN INIT INFO
### END INIT INFO

USER="root" # User we wil run application as.
CMD="sitestatus"
WORKDIR="/usr/local/site_installer"
LOGDIR="/var/log/site_installer"
PIDDIR="/var/run/site_installer"
NAME="site_installerd" 

###### Start script ########################################################

recursiveKill() { # Recursively kill a process and all subprocesses
    CPIDS=$(pgrep -P $1);
    for PID in $CPIDS
    do
        recursiveKill $PID
    done
    sleep 3 && kill -9 $1 2>/dev/null & # hard kill after 3 seconds
    kill $1 2>/dev/null # try soft kill first
}

case "$1" in
      start)
        echo "Starting $NAME ..."
        if [ -f "$PIDDIR/$NAME.pid" ]
        then
            echo "Already running according to $PIDDIR/$NAME.pid"
            exit 1
        fi
        cd "$WORKDIR"
        /bin/su -m -l $USER -c "$CMD" > "$LOGDIR/$NAME.log" 2>&1 &
        PID=$!
        echo $PID > "$PIDDIR/$NAME.pid"
        echo "Started with pid $PID - Logging to $LOGDIR/$NAME.log" && exit 0
        ;;
      stop)
        echo "Stopping $NAME ..."
        if [ ! -f "$PIDDIR/$NAME.pid" ]
        then
            echo "Already stopped!"
            exit 1
        fi
        PID=`cat "$PIDDIR/$NAME.pid"`
        recursiveKill $PID
        rm -f "$PIDDIR/$NAME.pid"
        echo "stopped $NAME" && exit 0
        ;;
      restart)
        $0 stop
        sleep 1
        $0 start
        ;;
      status)
        if [ -f "$PIDDIR/$NAME.pid" ]
        then
            PID=`cat "$PIDDIR/$NAME.pid"`
            if [ "$(/bin/ps --no-headers -p $PID)" ]
            then
                echo "$NAME is running (pid : $PID)" && exit 0
            else
                echo "Pid $PID found in $PIDDIR/$NAME.pid, but not running." && exit 1
            fi
        else
            echo "$NAME is NOT running" && exit 1
        fi
    ;;
      *)
      echo "Usage: /etc/init.d/$NAME {start|stop|restart|status}" && exit 1
      ;;
esac

exit 0