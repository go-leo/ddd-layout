#!/bin/sh
function echo_name()
{
  a=()
	#shell会执行反引号中的内容（命令）
	for file in $1
	do
	  [[ -e "$file" ]] || break  # handle the case of no *.wav files
    echo "$file"
		if [ -d $1'/'$file ]
		then
#		    a[${#a[*]}] = $1'/'$file
			echo_name $1'/'$file
		else
			echo $1'/'$file
		fi
	done
#	echo ${a[*]}
}
echo_name $1