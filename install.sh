#!/bin/bash
clear
echo '-------------------------------'
echo '|  Configuring Easy-V2ray-Mu  |'
echo '-------------------------------'
echo 'Please enter your Node ID:'
read node_id
echo 'Please enter your Mu-api URI(eg:http://www.xxxx.com/mu/v2):'
read mu_uri
echo 'Please enter your Mu-api KEY:'
read mu_key
echo '-------------------------------'
echo '|        Your Configure       |'
echo '-------------------------------'
echo 'Your Node ID:'
echo $node_id
echo 'Your Mu-api URI:'
echo $mu_uri
echo 'Your Mu-api KEY:'
echo $mu_key
echo 'Is it OK?(y/n)'
read isok
if [ $isok != 'y' -a $isok != 'Y' ];
then 
	echo 'Quit Install'
	exit
fi
echo '-------------------------------'
echo '|        Installing...        |'
echo '-------------------------------'
yum install wget
wget https://github.com/v2ray/v2ray-core/releases/download/v3.27/v2ray-linux-64.zip
unzip v2ray-linux-64.zip
rm -rf v2ray-linux-64.zip
mv v2ray-v3.27-linux-64 v2ray-mu
cd v2ray-mu
mkdir log
touch log/error.log
touch log/access.log
touch log/v2ray-mu.log
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/cfg.json
wget https://github.com/tonychanczm/easy-v2ray-mu/releases/download/v1.1/v2mctl
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/mu.conf
sed -i "s;##mu_uri##;$mu_uri;g" mu.conf
sed -i "s;##mu_key##;$mu_key;g" mu.conf
sed -i "s;##node_id##;$node_id;g" mu.conf


