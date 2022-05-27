#! /bin/bash
set -e

if [ $1 == 5 ]
then
  ./paxosserver -laddr="localhost:50081" -addrs="localhost:50083,localhost:50082,localhost:50084,localhost:50085" &
  ./paxosserver -laddr="localhost:50082" -addrs="localhost:50083,localhost:50081,localhost:50084,localhost:50085" &
  ./paxosserver -laddr="localhost:50083" -addrs="localhost:50081,localhost:50082,localhost:50084,localhost:50085" &
  ./paxosserver -laddr="localhost:50084" -addrs="localhost:50083,localhost:50082,localhost:50081,localhost:50085" &
  ./paxosserver -laddr="localhost:50085" -addrs="localhost:50083,localhost:50082,localhost:50084,localhost:50081" &
else
  ./paxosserver -laddr="localhost:50081" -addrs="localhost:50083,localhost:50082" &
  ./paxosserver -laddr="localhost:50082" -addrs="localhost:50081,localhost:50083" &
  ./paxosserver -laddr="localhost:50083" -addrs="localhost:50081,localhost:50082" &
fi

echo "running, enter to stop"

read && killall paxosserver
