#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "~/github.com/sbanka1/lenslocked"

echo "==== Releasing lenslocked ===="
echo " Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslocked
echo " Done!"

echo " Deleting existing code..."
ssh root@sbanka "rm -rf /root/lenslocked"
echo " Code deleted successfully!"

echo " Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@sbanka:/root/lenslocked/
echo " Code uploaded successfully!"

echo " Building the code on remote server..."
ssh root@sbanka 'export GOPATH=/root/go; cd /root/lenslocked; /usr/local/go/bin/go build -o /root/app/server *.go'
echo " Code built successfully!"

echo " Moving assets..."
ssh root@sbanka "cd /root/app; cp -R /root/lenslocked/assets ."
echo " Assets moved successfully!"

echo " Moving views..."
ssh root@sbanka "cd /root/app; cp -R /root/lenslocked/views ."
echo " Views moved successfully!"

echo " Restarting the server..."
ssh root@sbanka "sudo service lenslocked restart"
echo " Server restarted successfully!"

echo " Restarting Caddy server..."
ssh root@sbanka "sudo service caddy restart"
echo " Caddy restarted successfully!"

echo "==== Done releasing lenslocked ===="

