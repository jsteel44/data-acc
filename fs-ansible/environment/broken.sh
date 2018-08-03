#!/bin/bash

. .venv/bin/activate

set -eux

echo Get FS1 up

ansible-playbook test-dac.yml -i test-inventory --tag format_mgs --tag reformat_mdts --tag reformat_osts
ansible-playbook test-dac.yml -i test-inventory --tag start_mgs --tag start_mdts --tag start_osts --tag mount_fs

ls /mnt/lustre/fs1
ssh dac2 'sudo bash -c "hostname > /mnt/lustre/fs1/demo"'
cat /mnt/lustre/fs1/demo

echo Get FS2 up

ansible-playbook test-dac2.yml -i test-inventory2 --tag format_mgs --tag reformat_mdts --tag reformat_osts
ansible-playbook test-dac2.yml -i test-inventory2 --tag start_mgs --tag start_mdts --tag start_osts --tag mount_fs

ls /mnt/lustre/fs2
ssh dac2 'sudo bash -c "hostname > /mnt/lustre/fs2/demo"'
cat /mnt/lustre/fs2/demo


echo Get FS2 down

ansible-playbook test-dac2.yml -i test-inventory2 --tag umount_fs --tag stop_osts --tag stop_mdts
ansible-playbook test-dac2.yml -i test-inventory2 --tag reformat_mdts --tag reformat_osts

cat /mnt/lustre/fs2/demo || true

echo Get FS1 down

ansible-playbook test-dac.yml -i test-inventory --tag umount_fs --tag stop_osts --tag stop_mdts
ansible-playbook test-dac.yml -i test-inventory --tag reformat_mdts --tag reformat_osts


cat /mnt/lustre/fs1/demo || true


echo Get FS2 up

ansible-playbook test-dac2.yml -i test-inventory2 --tag format_mgs --tag reformat_mdts --tag reformat_osts
ansible-playbook test-dac2.yml -i test-inventory2 --tag start_mgs --tag start_mdts --tag start_osts --tag mount_fs

ls /mnt/lustre/fs2
ssh dac2 'sudo bash -c "hostname > /mnt/lustre/fs2/demo"'
cat /mnt/lustre/fs2/demo


echo Get FS2 down

ansible-playbook test-dac2.yml -i test-inventory2 --tag umount_fs --tag stop_osts --tag stop_mdts
ansible-playbook test-dac2.yml -i test-inventory2 --tag reformat_mdts --tag reformat_osts

cat /mnt/lustre/fs2/demo || true
