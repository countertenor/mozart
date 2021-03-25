#!/bin/bash	

echo "{{.values.value1}} {{.values.value2}}"	
pwd	
ls

echo "{{.file3}}" > /tmp/a.sh
bash /tmp/a.sh