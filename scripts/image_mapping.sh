# the ip address depends on docker image property
sudo sh -c 'echo "172.17.0.3 my-mysql" >> /etc/hosts'
sudo sh -c 'echo "172.17.0.4 my-redis" >> /etc/hosts'