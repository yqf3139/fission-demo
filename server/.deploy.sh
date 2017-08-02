#!/usr/bin/env bash

fission env create --name python --image yqf3139/python3-env
fission fn create --name app-client-filter --code ./app-client-filter.js --env node
fission fn create --name app-image-recognize --code ./app-image-recognize.py --env python
