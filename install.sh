#!/bin/bash
#默认通讯端口为8300，协议为ws
#下载自定面板配置文件，文件格式请见本目录下的config.conf
#wget -O config.conf https://www.baidu.com/config.conf
echo '*****NOTICE*****'
echo 'Please **EDIT THE CONFIG FILE** (config.conf) Before You Install'
cd v2ray
chmod +x install.sh
./install.sh
echo 'echo "You have been install, nothing to do"' > install.sh
cd ..
echo 'Install Success! Please type "sh restart.sh" to run v2ray-mu'
