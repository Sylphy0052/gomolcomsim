#!/bin/bash

function setStarttime() {
        start_time=`date +%s`
}

function getEndtime() {
        end_time=`date +%s`

        SS=`expr ${end_time} - ${start_time}`

        HH=`expr ${SS} / 3600`
        SS=`expr ${SS} % 3600`
        MM=`expr ${SS} / 60`
        SS=`expr ${SS} % 60`

        echo "${HH}:${MM}:${SS}" >> ./result/time.txt
        echo "${HH}:${MM}:${SS}"
}

read -p "Count(default=1000): " input
expr $input + 1 > /dev/null 2>&1
RET=$?
if [ $RET -lt 2 ]; then
    count=$input
else
    count=1000
fi

read -p "Wait (y/N): " yn
if [[ $yn = 'y' ]] ; then
  isptime="true"
else
  isptime="false"
fi

num=`ls -U1 ./dat/*.dat | wc -l`
file_num=`expr $num`
file_index=0

mkdir -p result

touch result/time.txt
echo '' > ./result/time.txt

go build execute.go

for file in ./dat/*.dat
do
	echo "./execute $file $isptime"
	file_index=$(( file_index + 1 ))
  SECONDS=0
  echo "$file" >> ./result/time.txt
	for i in `seq 1 $count`
	do
    setStarttime
		echo "Count: $i / $count --- $file_index of $file_num"
		./execute $file $isptime
    getEndtime
	done
  echo "$file: $SECONDS / $count" >> ./result/time.txt
  echo "$file: $SECONDS / $count"
done
