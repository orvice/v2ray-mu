#!/bin/bash
cd log
rm -rf access.log
touch access.log
rm -rf error.log
touch error.log
rm -rf v2ray-mu.log
touch v2ray-mu.log
cd ..
echo "All Logs Clear!"
