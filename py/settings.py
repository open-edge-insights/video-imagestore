from collections import namedtuple

defaults = {
    'host': 'localhost',
    'port': '6379',
    'retentionpolicy': 400,
    'inmemory': 'redis'
}

value = namedtuple('inmemory', defaults.keys())(**defaults)
