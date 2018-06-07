#!/bin/bash
echo 'Please enter your Node ID:'
read node_id
echo 'Install Start!'
yum -y install docker-io
systemctl enable docker
systemctl start docker
yum -y install epel-release
yun -y install python
yum -y install python-pip
pip install --upgrade pip
pip --default-timeout=200 install -U docker-compose
chmod -R 777 log/
chmod +x run.sh
chmod +x config.sh
chmod +x stop.sh
echo "/root/v2ray/v2ray/run.sh" >> /etc/rc.d/rc.local
chmod +x /etc/rc.d/rc.local
echo 'Configing...'
sed -i "s/#ID_OF_NODE#/$node_id/g" docker-compose.yml
sed -i "s/#ID_OF_NODE#/$node_id/g" mu.env
./config.sh
