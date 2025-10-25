#!bin/bash
cd ../HomeAssistant/
export DOCKER_HOST=unix:///var/run/docker.sock
/usr/bin/docker compose up -d

cd ../ telegrammBot/
/usr/bin/docker compose up -d --build