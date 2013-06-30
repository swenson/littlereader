#!/bin/bash

rsync -avz --delete requirements.sh Makefile src swenson@swenson.caswenson.com:/u/apps/rss/ &&
ssh swenson@swenson.caswenson.com '/bin/bash --login -c "cd /u/apps/rss; make reader"' &&
ssh root@swenson.caswenson.com '/etc/init.d/rss restart'