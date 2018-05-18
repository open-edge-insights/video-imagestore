from collections import namedtuple

defaults = {
    'host': '10.223.97.21',
    'port': '6379',
    'retentionpolicy': 400,
    'inmemory': 'redis'
}

value = namedtuple('inmemory', defaults.keys())(**defaults)
