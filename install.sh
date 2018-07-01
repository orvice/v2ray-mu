#!/bin/bash
#a strange script to install v2ray-mu
clear
mu_uri=$1
mu_key=$2
node_id=$3
echo '-------------------------------'
echo '|  Configuring Easy-V2ray-Mu  |'
echo '-------------------------------'
if [ ! $node_id ];
then 
	echo 'Please enter your Node ID:'
	read node_id
fi
if [ ! $mu_uri ];
then 
	echo 'Please enter your Mu-api URI(eg:http://www.xxxx.com/mu/v2):'
	read mu_uri
fi
if [ ! $mu_key ];
then 
	echo 'Please enter your Mu-api KEY:'
	read mu_key
fi
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
yum install unzip -y
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
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/run.sh
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/stop.sh
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/cleanLogs.sh
wget https://raw.githubusercontent.com/tonychanczm/easy-v2ray-mu/dev/catLogs.sh
chmod +x *
echo "30 4 * * * cd $(readlink -f .) && ./run.sh">> /var/spool/cron/root
echo '--------------------------------'
echo '|       Install finshed        |'
echo '|please run this command to run|'
echo '----------- V  V  V ------------'
echo "cd $(readlink -f .) && ./run.sh"

