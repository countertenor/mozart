#!/bin/bash	

for i in {1..5}	
do	
    echo "printing $i"	
    sleep 1
done

DIR={{.get_current_dir}}
bash $DIR/!external_file.sh