#! /usr/bin/env python
# -*- coding: utf-8 -*-
#
# Copyright © 2019 Yongwen Zhuang <zeoman@163.com>
#
# Distributed under terms of the MIT license.

"""
S10

"""
import requests

url = "https://api.live.bilibili.com/guard/topList"
params = dict(
    roomid=5085,
    page=1,
    ruid=137952
)
resp = requests.get(url=url, params=params)
data = resp.json()
pages = data['data']['info']['page']

top3 = data['data']['top3']

for page in range(pages):
    print(page+1)
