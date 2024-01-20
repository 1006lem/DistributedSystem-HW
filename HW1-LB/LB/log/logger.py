# logger.py
# : Define Logging format (time, module, function info)
# ----------------------------------------------------------

import logging


class Setting:
    LEVEL = logging.INFO
    FORMAT = '[%(levelname)s] %(asctime)s  "%(message)s"  [MODULE %(module)s] in  %(filename)s #%(lineno)d : func %(funcName)s(...)'


def Logger():
    logger = logging.getLogger('info_logger')
    logger.setLevel(Setting.LEVEL)

    # Logging format
    formatter = logging.Formatter(Setting.FORMAT)
    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(formatter)
    logger.addHandler(stream_handler)
    return logger


def info(message):
    Logger().info("%s" % (str(message)))


def warning(message):
    Logger().warning("%s" % (str(message)))


def error(message):
    Logger().error("%s" % (str(message)))
