#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/github.com/sbanka1/lenslocked"

echo "==== Releasing lenslocked ===="
echo " Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslocked
echo " Done!"

echo " Deleting existing code..."
ssh root@lenslocked.sbanka.io "rm -rf /root/go/src/github.com/sbanka1/lenslocked"
echo " Code deleted successfully!"

echo " Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@lenslocked.sbanka.io:/root/go/src/github.com/sbanka1/lenslocked/
echo " Code uploaded successfully!"

echo " Go getting deps..."
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@lenslocked.sbanka.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

echo " Building the code on remote server..."
ssh root@lenslocked.sbanka.io 'export GOPATH=/root/go; cd /root/app; /usr/local/go/bin/go build -o ./server $GOPATH/src/github.com/sbanka1/lenslocked/*.go'
echo " Code built successfully!"

echo " Moving assets..."
ssh root@lenslocked.sbanka.io "cd /root/app; cp -R /root/go/src/github.com/sbanka1/lenslocked/assets ."
echo " Assets moved successfully!"

echo " Moving views..."
ssh root@lenslocked.sbanka.io "cd /root/app; cp -R /root/go/src/github.com/sbanka1/lenslocked/views ."
echo " Views moved successfully!"

echo " Moving Caddyfile..."
ssh root@lenslocked.sbanka.io "cd /root/app; cp /root/go/src/github.com/sbanka1/lenslocked/Caddyfile ."
echo " Caddyfile moved successfully!"

echo " Restarting the server..."
ssh root@lenslocked.sbanka.io "sudo service lenslocked restart"
echo " Server restarted successfully!"

echo " Restarting Caddy server..."
ssh root@lenslocked.sbanka.io "sudo service caddy restart"
echo " Caddy restarted successfully!"

echo "==== Done releasing lenslocked ===="

