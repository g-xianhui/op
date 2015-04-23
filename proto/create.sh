#!/bin/sh

# server
protoc --go_out=../server/pb/ *.proto

# client
protoc --go_out=../client/pb/ *.proto

file="msgtype.go"
rm -f $file
touch $file
echo "package pb" >> $file
echo "const (" >> $file
echo "    _ = iota" >> $file

cat "msgtype.txt" | while read line
do
    echo "    $line" >> $file
done

echo ")" >> $file

cp -f $file ../server/pb
cp -f $file ../client/pb

rm -f $file
