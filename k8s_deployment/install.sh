#!/bin/bash

ansible-playbook  install_plb.yaml  -e dsm_username=$DSM_USERNAME -e dsm_password=$DSM_PASSWORD -e image_pull_secret=$PULL_SECRET