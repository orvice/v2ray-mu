#!/bin/bash
source ../config.conf
sed -i "s;#PANEL_ADDR#;$panel_url;g" docker-compose.yml
sed -i "s;#PANEL_ADDR#;$panel_url;g" mu.env
sed -i "s/#MU_KEY#/$mu_key/g" docker-compose.yml
sed -i "s/#MU_KEY#/$mu_key/g" mu.env
