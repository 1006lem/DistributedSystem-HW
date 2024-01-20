import os

from LB.control.control import start_server
from LB.log.logger import Logger

Log = Logger()
if __name__ == "__main__":
    Log.info("Start LB server")
    print()
    print("-------------------------------------------------------------------------")
    print(" _                        _  ______  _")
    print("| |                      | | | ___ \| |")
    print("| |      ___    __ _   __| | | |_/ /| |  __ _  _ __    ___   ___  _ __")
    print("| |     / _ \  / _` | / _` | | ___ \| | / _` || '_ \  / __| / _ \| '__|")
    print("| |____| (_) || (_| || (_| | | |_/ /| || (_| || | | || (__ |  __/| |")
    print("\_____/ \___/  \__,_| \__,_| \____/ |_| \__,_||_| |_| \___| \___||_|")
    print("-------------------------------------------------------------------------")
    print()
    print()
    control_server_port = int(os.environ.get('CONTROL_SERVER_PORT', "8080"))
    start_server(control_server_port)
