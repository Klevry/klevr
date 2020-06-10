#!/bin/bash
### docker image build automation with verification
### Made by ralf.yang@gsshop.com, goody80762@gmail.com
### Version 0.5


## YOU CAN CHANGE THE COMMAND for WHAT YOU NEED AS BELOW CheckCommand.
#Docker_repo=`cat /data/z/etc/init.d/docker | grep "^connection" | awk -F '=' '{print $2}' | sed -e 's/"//g' `
Docker_repo=""
CheckCommand="ls -lai"
VerifiedTag="VF"

Output="/tmp/build_auto.list"
BAR="====================================="


WorkDir="$PWD"
OutputDir="$WorkDir/Output"
	if [[ `(which rpm)` != "" ]];then
		sudo_pkg_check=`rpm -qa  2> /dev/null |grep "^sudo"`
		InstM="rpm -qa |grep 'vim-'"
	elif [[ `(which apt-get)` != "" ]];then
		sudo_pkg_check=`dpkg -l  2> /dev/null |grep "^sudo"`
		InstM="dpkg -l |grep 'vim-'"
	fi

## Please modify under variable for the private repository if you have.
	if [[ $sudo_pkg_check = "" ]]; then
		Comm="docker"
	else
		Comm="sudo docker"
	fi


ImagesList=(`ls -d */ |egrep -v "Output" | sed -e 's#/##g'`)
	if [[ $ImagesList = "" ]];then
		echo ""
		echo " === Please make sure the target Dockerfile or directory ! ==="
		echo ""
		exit 0;
	fi

clear
	Count=0
	echo "" > $Output
	echo "$BAR" >> $Output
	while [ $Count -lt ${#ImagesList[@]} ]; do
		echo "$((Count+1)) | ${ImagesList[$Count]}" >> $Output
	let Count=$Count+1
	done
	echo "$BAR" >> $Output
## Menu
cat $Output
echo "Please insert a number for build:"
read build_num

## Select a number for work
DockerName=`grep "^$build_num" $Output | awk '{print $3}'`

## Checks for changes
cp ./$DockerName/Dockerfile ./$DockerName/.Dockerfile_bak

	if [[ $InstM != "" ]];then
		ViComm=`which vim | tail -1`
	else
		ViComm=`which vi | tail -1`
	fi

$ViComm ./$DockerName/Dockerfile
	if [[ `(diff ./$DockerName/Dockerfile ./$DockerName/.Dockerfile_bak)` = "" ]];then
		echo "$BAR"
		echo " Nothing has been changed!! processor will be stop"
		rm -f ./$DockerName/.Dockerfile_bak
		read
		exit 0
	fi
rm -f ./$DockerName/.Dockerfile_bak


echo "$BAR"
echo " If you don't go anymore, plese insert a [ n ] key or break"
read keystop
	if [[ $keystop = "n" ]];then
		echo "$BAR"
		echo " Build processor has been stopped !"
		read
		exit 0
	fi


##
mkdir $OutputDir -p 2> /dev/null

DockerNameOut="$OutputDir/$DockerName.out"
DockerUplog="$OutputDir/$DockerName.log"
ErrorOut="buildauto.err"
EnablePortArry=(`cat ./$DockerName/Dockerfile |grep EXPOSE | sed -e 's/EXPOSE //g'`)

## Find out the version by 'LABEL version=x.x'
Tags=`grep " version=" ./$DockerName/Dockerfile | awk -F '=' '{print $2}' | sed 's/"//g'`


## Check build success
$Comm build -t $DockerName:$Tags ./$DockerName |tee $DockerNameOut
	if [[ `(grep  "Successfully built" $DockerNameOut)` = "" ]];then
		echo "$BAR"
		echo "Build failed!!!"
		echo "$BAR"
		exit 0;
	fi

## Running docker for check
$Comm run -d --name $DockerName $DockerName:$Tags tailf /etc/resolv.conf
	if [[ `($Comm exec $DockerName $CheckCommand)` = "" ]];then
		echo "$BAR"
		echo " Zinst pacakge dose not install in Docker image yet!!"
		echo "$BAR"
		exit 0;
	fi


	if [[ ${EnablePortArry[@]} != "" ]];then
		#countPort=0
		#while [ $countPort -lt ${#EnablePortArry[@]} ];do
			#$Comm exec -it $DockerName netstat -antp |grep ${EnablePortArry[$countPort]}
			if [[ `($Comm ps |grep ${EnablePortArry[0]})` = "" ]];then
				echo ""
				echo "$BAR"
				echo "EXPOSE Port has not attented"
				echo "$BAR"
				exit 0;
			else
				echo " === Port EXPOSE check==="
				$Comm ps |grep ${EnablePortArry[0]}
			fi
		#let countPort=$countPort+1
		#done
	fi


## Everything is fine
echo " Stopping the Temporary container....."
$Comm stop $DockerName
$Comm rm $DockerName

echo ""

ImgChk=`$Comm images |grep "^$VerifiedTag/$DockerName"`
	if [[ $ImgChk != "" ]];then
		echo "$BAR"
		echo "Image name already existed!"
		echo "$BAR"
		#$Comm rmi  -f "$VerifiedTag//$DockerName"
	fi
echo " Running new container for tagging new...."

$Comm tag $DockerName:$Tags $VerifiedTag/$DockerName:$Tags

echo "$BAR"
echo " Good job!  everything is okay!"
$Comm images |grep "^$VerifiedTag/$DockerName"
echo "$BAR"

rm -f $DockerNameOut

$Comm tag $VerifiedTag/$DockerName:$Tags $Docker_repo/$VerifiedTag/$DockerName:$Tags
$Comm tag $Docker_repo/$VerifiedTag/$DockerName:$Tags $Docker_repo/$VerifiedTag/$DockerName:latest
$Comm push $Docker_repo/$VerifiedTag/$DockerName:$Tags |tee $ErrorOut
$Comm push $Docker_repo/$VerifiedTag/$DockerName:latest |tee $ErrorOut


## Check upload
        if [[ `(grep  "connection refused" $ErrorOut)` != "" ]];then
                echo "$BAR"
                echo " Upload failed!!! You should check a Registry."
                echo "$BAR"
		$Comm rmi -f $Docker_repo/$VerifiedTag/$DockerName:$Tags
		$Comm rmi -f $VerifiedTag/$DockerName:$Tags
		$Comm rmi -f $Docker_repo/$VerifiedTag/$DockerName:latest
                exit 0;
        fi
## Upload log delete
rm -fR $ErrorOut

### Parsing for IP added address
#curl http://$Docker_repo/v2/_catalog 2> /dev/null | sed -e 's/[{}]/''/g' | awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' | sed -e 's#\[##g' -e 's#\]##g' -e 's#"repositories"\:##g' -e "s#^\"#$Docker_repo/#g" -e 's/"//g' |sort --version-sort |grep "/gs/"

	ImgArry=(`curl http://$Docker_repo/v2/_catalog 2> /dev/null | sed -e 's/[{}]/''/g' | awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' | sed -e 's#\[##g' -e 's#\]##g' -e 's#"repositories"\:##g'  -e 's/"//g' | grep "$VerifiedTag/"`)
	Counter=0
	TempArry=./temp
	ResultList=$WorkDir/available_image.list
	touch $TempArry
	while [ $Counter -lt ${#ImgArry[@]} ];do
		curl "http://$Docker_repo/v2/${ImgArry[$Counter]}/tags/list" 2> /dev/null | sed -e 's/[{}]/''/g' | awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' |sed -e 's/"//g' >> $TempArry
		cat $TempArry 2> /dev/null | sed -e 's/[{}]/''/g' | awk -v k="text" '{n=split($0,a,","); for (i=1; i<=n; i++) print a[i]}' | sed -e 's/"//g' -e 's/name:/\n---/g' -e 's/tags://g' -e 's/\[//g' -e 's/\]//g' > $ResultList
	let Counter=$Counter+1
	done
rm -Rf $TempArry $OutputDir 2> /dev/null

$Comm rmi -f $Docker_repo/$VerifiedTag/$DockerName:$Tags
$Comm rmi -f $VerifiedTag/$DockerName:$Tags
$Comm rmi -f $Docker_repo/$VerifiedTag/$DockerName:latest

cd $WorkDir/../; git pull;git add *;git commit -m "update the Available list for - $Docker_repo/$VerifiedTag/$DockerName:$Tags" ; git push


