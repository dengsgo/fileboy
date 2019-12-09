#!/bin/bash

# 合并上游-fork来源
git remote add upstream https://github.com/dengsgo/fileboy.git
git fetch upstream
git checkout master
git merge upstream/master