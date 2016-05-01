#!/usr/bin/python3

import os.path as path
import socket
import subprocess as sp
import psutil

NODE = "/usr/bin/node"
SERVER = "/opt/cloud9/c9sdk/server.js"
USER = sp.check_output(["whoami"]).decode().rstrip()
HOME = path.expanduser("~")


def is_cloud9(pinfo):
    cmd = pinfo["cmdline"]
    return pinfo["username"] == USER and\
        cmd[0] == NODE and\
        cmd[1] == SERVER and\
        cmd[2] == "-p"


def get_cloud9():
    pid, port = None, None
    for proc in psutil.process_iter():
        try:
            pinfo = proc.as_dict(attrs=["pid", "username", "cmdline"])
        except psutil.NoSuchProcess:
            pass
        else:
            if is_cloud9(pinfo):
                pid, port = pinfo["pid"], int(pinfo["cmdline"][3])
    return pid, port


def free_port():
    s = socket.socket()
    s.bind(("", 0))
    port = s.getsockname()[1]
    s.close()
    return port


def start_cloud9(port=None):
    if port is None:
        port = free_port()
    with open(path.join(HOME, ".cloud9.log"), "w") as fh:
        cmd = [NODE, SERVER, "-p", str(port), "-w", HOME]
        proc = sp.Popen(cmd, stdout=fh, stderr=fh)
    return proc.pid, port


if __name__ == "__main__":
    pid, port = get_cloud9()
    if pid is None or port is None:
        pid, port = start_cloud9()
    print(pid, port, sep="\n")
