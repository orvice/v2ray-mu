#!/bin/bash
v2ray_pid=$(ps ux | grep "/root/v2ray/v2ray/v2ray" | grep -v grep | awk '{print $2}')
v2ray_mu_cid=$(docker ps -a | grep "./v2ray-mu" | grep -v grep | awk '{print $1}')

if [ ! $v2ray_pid ];
then
    echo 'V2ray is not running, nothing to do.'
else
    echo 'Stopping V2Ray (pid:'$v2ray_pid')'
    kill -9 $v2ray_pid
    echo 'V2Ray is DOWN'
fi

if [ ! $v2ray_mu_cid ];
then
echo 'V2Ray-mu Manager is not running, nothing to do.'
else
echo 'Stopping V2Ray-mu Manager (cid:'$v2ray_mu_cid')'
docker stop $v2ray_mu_cid
docker rm $v2ray_mu_cid
echo 'V2Ray-mu Manager is DOWN'
fi
