#!/bin/bash

# 须以特权模式运行（docker run --privileged）
sysctl net.core.somaxconn=1024

/usr/local/bin/ddgo $@