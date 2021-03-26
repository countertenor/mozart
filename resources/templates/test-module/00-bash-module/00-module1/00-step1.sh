#!/bin/bash	

echo "{{.values.value1}} {{.values.value2}}"	
pwd	
ls

{{writeFile "/tmp/a.sh" .file3}}