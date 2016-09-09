#!/bin/bash

echo "mode: set" > acc.cov
for Dir in $(find ./* -maxdepth 10 -type d ); 
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -covermode=count -coverprofile=profile.cov $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.cov ]
    		then
        		cat profile.cov | grep -v "mode: set" >> acc.cov 
    		fi
    	else
    		exit 1
    	fi	
    fi
done
if [ -n "$COVERALLS" ]
then
	goveralls -coverprofile=acc.cov -service=travis-ci
fi	

# rm -rf ./profile.cov
# rm -rf ./acc.cov
