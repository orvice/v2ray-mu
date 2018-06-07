#!/bin/bash
touch /etc/systemd/system/v2ray-mu.service
cmd="sh /root/v2ray/v2ray/r.sh"
echo "[Unit]">/etc/systemd/system/v2ray-mu.service
echo "Description=V2ray deamon">>/etc/systemd/system/v2ray-mu.service
echo "[Service]">>/etc/systemd/system/v2ray-mu.service
echo "Type=simple">>/etc/systemd/system/v2ray-mu.service
echo "ExecStart=${cmd}">>/etc/systemd/system/v2ray-mu.service
echo "[Install]">>/etc/systemd/system/v2ray-mu.service
echo "WantedBy=multi-user.target">>/etc/systemd/system/v2ray-mu.service

