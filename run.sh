#!/bin/bash
v2ray_realpath=$(readlink -f v2ray)
cfg_realpath=$(readlink -f cfg.json)
v2muctl_realpath=$(readlink -f v2mctl)
v2ray_pid=$(ps ux | grep "$(readlink -f v2ray)" | grep -v grep | awk '{print $2}')
v2muctl_pid=$(ps ux | grep "$(readlink -f v2mctl)" | grep -v grep | awk '{print $2}')
if [ ! $v2ray_pid ];
then
    echo 'Starting V2Ray'
else
    echo 'Restarting V2Ray (pid:'$v2ray_pid')'
    kill -9 $v2ray_pid
fi

if [ ! $v2muctl_pid ];
then
echo 'Starting V2Ray-mu Manager'
else
echo 'Retarting V2Ray-mu Manager (pid:'$v2muctl_pid')'
kill -9 $v2muctl_pid
fi

cd log
rm -rf access.log
touch access.log
rm -rf error.log
touch error.log
rm -rf v2ray-mu.log
touch v2ray-mu.log
cd ..
echo "All Logs Clear!"

source mu.conf
export MU_URI=$MU_URI
export MU_TOKEN=$MU_TOKEN
export MU_NODE_ID=$NodeId
export SYNC_TIME=$SYNC_TIME
export V2RAY_ADDR=$V2RAY_ADDR
export V2RAY_TAG=$V2RAY_TAG



nohup $(readlink -f v2ray) --config=$(readlink -f cfg.json)>> /dev/null 2>&1 &
nohup $(readlink -f v2mctl)>> /dev/null 2>&1 &

sleep 1

v2ray_pid=$(ps ux | grep "$(readlink -f v2ray)" | grep -v grep | awk '{print $2}')
v2muctl_pid=$(ps ux | grep "$(readlink -f v2mctl)" | grep -v grep | awk '{print $2}')

if [ ! $v2ray_pid ];
then
    echo '***Fail to start V2Ray***'
else
    echo 'Success to start V2Ray (pid:'$v2ray_pid')'
fi

if [ ! $v2muctl_pid ];
then
echo '***Fail to start V2Ray-mu Manager***'
else
echo 'Success to start V2Ray-mu Manager (pid:'$v2muctl_pid')'
fi
