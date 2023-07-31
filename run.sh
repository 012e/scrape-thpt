#!/bin/bash
# begin: 51000001
# end: 51019942
echo running
HTTP_PROXY="https://103.141.143.102:41516" go run . -end=51000042 -start=51000001 -con=1
# HTTP_PROXY="https://103.141.143.102:41516" go run . -start=1 -end=11 -con=5