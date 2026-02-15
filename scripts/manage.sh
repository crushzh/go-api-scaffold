#!/bin/sh
#
# Service management script (Linux)
#
# Commands: start | stop | restart | status
# Usage: ./manage.sh {start|stop|restart|status}
#
# Auto-start (systemd):
#   cp myapp.service /etc/systemd/system/
#   systemctl enable myapp
#
# Auto-start (rcS):
#   Add "/opt/myapp/manage.sh start &" to the end of /etc/init.d/rcS
#

# ==================== Configuration ====================
APP_NAME="myapp"
APP_DIR="$(cd "$(dirname "$0")/.." && pwd)"
APP_BIN="${APP_DIR}/${APP_NAME}"
PID_FILE="${APP_DIR}/${APP_NAME}.pid"
WATCHDOG_PID_FILE="${APP_DIR}/${APP_NAME}-watchdog.pid"
WATCHDOG_LOG="${APP_DIR}/${APP_NAME}-watchdog.log"
CHECK_INTERVAL=60

# ==================== Helper functions ====================

log_msg() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

is_running() {
    local pid_file="$1"
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file" 2>/dev/null)
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
    fi
    return 1
}

get_pid() {
    local pid_file="$1"
    if [ -f "$pid_file" ]; then
        cat "$pid_file" 2>/dev/null
    fi
}

# ==================== Core operations ====================

do_start() {
    if is_running "$PID_FILE"; then
        log_msg "${APP_NAME} already running (PID: $(get_pid "$PID_FILE"))"
        return 0
    fi

    if [ ! -x "$APP_BIN" ]; then
        log_msg "error: ${APP_BIN} not found or not executable"
        return 1
    fi

    cd "$APP_DIR" || {
        log_msg "error: cannot enter directory ${APP_DIR}"
        return 1
    }

    ./"$APP_NAME" >/dev/null 2>&1 &
    local pid=$!
    echo "$pid" > "$PID_FILE"

    sleep 1
    if is_running "$PID_FILE"; then
        log_msg "${APP_NAME} started (PID: ${pid})"
    else
        log_msg "warning: ${APP_NAME} exited immediately after start, check logs"
        rm -f "$PID_FILE"
        return 1
    fi

    do_watchdog_start
    return 0
}

do_stop() {
    do_watchdog_stop

    if ! is_running "$PID_FILE"; then
        log_msg "${APP_NAME} not running"
        rm -f "$PID_FILE"
        return 0
    fi

    local pid=$(get_pid "$PID_FILE")
    log_msg "stopping ${APP_NAME} (PID: ${pid}) ..."

    kill "$pid" 2>/dev/null
    local count=0
    while [ $count -lt 10 ]; do
        if ! kill -0 "$pid" 2>/dev/null; then
            break
        fi
        sleep 1
        count=$((count + 1))
    done

    if kill -0 "$pid" 2>/dev/null; then
        log_msg "graceful shutdown timeout, force killing ..."
        kill -9 "$pid" 2>/dev/null
        sleep 1
    fi

    rm -f "$PID_FILE"
    log_msg "${APP_NAME} stopped"
    return 0
}

do_restart() {
    log_msg "restarting ${APP_NAME} ..."
    do_stop
    sleep 1
    do_start
}

do_status() {
    if is_running "$PID_FILE"; then
        log_msg "${APP_NAME} running (PID: $(get_pid "$PID_FILE"))"
    else
        log_msg "${APP_NAME} not running"
    fi

    if is_running "$WATCHDOG_PID_FILE"; then
        log_msg "watchdog running (PID: $(get_pid "$WATCHDOG_PID_FILE"))"
    else
        log_msg "watchdog not running"
    fi
}

# ==================== Watchdog ====================

do_watchdog_start() {
    if is_running "$WATCHDOG_PID_FILE"; then
        return 0
    fi

    _watchdog_loop >> "$WATCHDOG_LOG" 2>&1 &
    echo $! > "$WATCHDOG_PID_FILE"
    log_msg "watchdog started (PID: $!, interval: ${CHECK_INTERVAL}s)"
}

do_watchdog_stop() {
    if is_running "$WATCHDOG_PID_FILE"; then
        local pid=$(get_pid "$WATCHDOG_PID_FILE")
        kill "$pid" 2>/dev/null
        sleep 1
        kill -9 "$pid" 2>/dev/null
        rm -f "$WATCHDOG_PID_FILE"
        log_msg "watchdog stopped"
    fi
}

_watchdog_loop() {
    log_msg "watchdog loop started, checking every ${CHECK_INTERVAL} seconds"

    while true; do
        sleep "$CHECK_INTERVAL"
        if ! is_running "$PID_FILE"; then
            log_msg "[Watchdog] ${APP_NAME} not running, auto-restarting ..."
            cd "$APP_DIR" || continue
            ./"$APP_NAME" >/dev/null 2>&1 &
            local pid=$!
            echo "$pid" > "$PID_FILE"
            sleep 1
            if is_running "$PID_FILE"; then
                log_msg "[Watchdog] ${APP_NAME} restart success (PID: ${pid})"
            else
                log_msg "[Watchdog] ${APP_NAME} restart failed, retrying in ${CHECK_INTERVAL}s"
                rm -f "$PID_FILE"
            fi
        fi
    done
}

# ==================== Entry point ====================

case "$1" in
    start)   do_start   ;;
    stop)    do_stop    ;;
    restart) do_restart ;;
    status)  do_status  ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

exit $?
