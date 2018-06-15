#!/bin/bash
v2ray_pid=$(ps ux | grep "/root/v2ray/v2ray/v2ray" | grep -v grep | awk '{print $2}')
v2ray_mu_cid=$(docker ps -a | grep "./v2ray-mu" | grep -v grep | awk '{print $1}')
if [ ! $v2ray_pid ];
then
    echo 'Starting V2Ray'
else
    echo 'Restarting V2Ray (pid:'$v2ray_pid')'
    kill -9 $v2ray_pid
fi

if [ ! $v2ray_mu_cid ];
then
echo 'Starting V2Ray-mu Manager'
else
echo 'Retarting V2Ray-mu Manager (cid:'$v2ray_mu_cid')'
docker stop $v2ray_mu_cid
docker rm $v2ray_mu_cid
fi

chmod +x clearLogs.sh
./cleanLogs.sh

docker-compose up -d
nohup /root/v2ray/v2ray/v2ray --config=/root/v2ray/v2ray/cfg.json>> /dev/null 2>&1 &

v2ray_pid=$(ps ux | grep "/root/v2ray/v2ray/v2ray" | grep -v grep | awk '{print $2}')
v2ray_mu_cid=$(docker ps -a | grep "./v2ray-mu" | grep -v grep | awk '{print $1}')

if [ ! $v2ray_pid ];
then
    echo '***Fail to start V2Ray***'
else
    echo 'Success to start V2Ray (pid:'$v2ray_pid')'
fi

if [ ! $v2ray_mu_cid ];
then
echo '***Fail to start V2Ray-mu Manager***'
else
echo 'Success to start V2Ray-mu Manager (cid:'$v2ray_mu_cid')'
fi
