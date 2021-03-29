#!/bin/bash	

echo "{{.values.value1}} {{.values.value2}}"	
pwd	
ls

{{writeContentToFile .file3 "mozart-generated/a.sh"}}

./mozart-generated/a.sh