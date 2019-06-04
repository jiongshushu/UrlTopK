#!/usr/bin/env python
# -*- coding: utf-8 -*-
########################################################################
#
# Copyright (c) 2019 Bytedance.com, Inc. All Rights Reserved
#
########################################################################

import random

hosts = [
    "www.baidu.com",
    "www.google.com",
    "www.hao123.com",
    "mail.google.com"
]

paths = [
    "/aaaa/",
    "/badasfdas/",
    "/dfdf/",
]

params = [
    "aa=bb",
    "zcasfd=2323&bb=11",
    "cccc=123&ddd=23213&ddd=888"
]

def main():

    f = open("test.txt", "a")

    url_num_map = {}
    for i in range(0, 11):
        host = hosts[random.randint(0, len(hosts)-1)]
        path = paths[random.randint(0, len(paths)-1)]
        query = params[random.randint(0, len(params)-1)]

        url = "{0}{1}?{2}\n".format(host, path, query)
        if url not in url_num_map:
            url_num_map[url] = 0
        url_num_map[url] += 1

        f.write(url)

    f.close()

    for k, v in url_num_map.items():
        print("%d %s" % (v, k))


if __name__ == "__main__":
    main()
    exit(0)
