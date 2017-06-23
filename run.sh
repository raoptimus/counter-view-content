#!/bin/sh
sudo mkdir /var/empty/run
sudo chown nobody:nobody /var/empty/run
sudo -H -u nobody  bash -c './bin/counter -config=./config/release.json'
