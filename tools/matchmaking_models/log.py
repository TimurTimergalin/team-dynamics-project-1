import logging
import logging.config
from io import StringIO

logs_stream = StringIO()


logging_config = {
    'version': 1,
    'disable_existing_loggers': False,
    'formatters': {
        'identity': {
            'format': '%(message)s'
        },
        'simple': {
            'format': '%(levelname)s: %(message)s'
        }
    },
    'handlers': {
        'stdout': {
            'level': 'DEBUG',
            'class': 'logging.StreamHandler',
            'formatter': 'simple',
            'stream': 'ext://sys.stdout'
        },
        'in_memory_log': {
            'level': 'INFO',
            'class': 'logging.StreamHandler',
            'formatter': 'identity',
            'stream': logs_stream
        }
    },
    'loggers': {
        'root': {
            'level': "DEBUG",
            'handlers': ['stdout', 'in_memory_log']
        }
    }
}

logger = logging.getLogger("in_memory_log")
logging.config.dictConfig(logging_config)


def get_in_memory_logs():
    return logs_stream.getvalue()


def clear_in_memory_logs():
    logs_stream.seek(0)
    logs_stream.truncate(0)
