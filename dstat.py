"""
Transform dstat JSON input into a format suitable for elastic search
"""

import copy
import json
import re
import sys


def transform(s):
    matches = re.findall(r'(\d+)(\.\d+)*(\S*)', s)
    number, decimals, unit = matches[0]

    def t2(number, decimals, mul=1):
        if decimals == '':
            return int(number) * mul
        else:
            return float(number + decimals) * mul

    if unit.upper() == '':
        return t2(number, decimals)
    elif unit.upper() == 'B':
        return t2(number, decimals)
    elif unit.upper() == 'K':
        return t2(number, decimals, mul=1024)
    elif unit.upper() == 'M':
        return t2(number, decimals, mul=1024 * 1024)
    elif unit.upper() == 'G':
        return t2(number, decimals, mul=1024 * 1024 * 1024)
    else:
        raise ValueError('Unable to convert {}'.format(s))


def main():
    d = json.loads(sys.stdin.read())
    d_copy = copy.deepcopy(d)

    def parse_item(item):
        if isinstance(item, str):
            res = transform(item)
            return res
        elif isinstance(item, dict):
            for k, v in dict(item).items():
                item[k] = parse_item(v)
            return item
        else:
            raise ValueError('Cannot handle type {}'.format(type(item)))

    parse_item(d_copy)
    sys.stdout.write(json.dumps(d_copy))


if __name__ == '__main__':
    main()
