#!/bin/bash
primary_cmd="$GOPATH/bin/vrgo --mode=server --port=1000 --backup_ports=9000 --backup_ports=9001 > primary.log&"

backup_cmd="$GOPATH/bin/vrgo --mode=backup"
backups=( "2","9000" "3","9001" )

client_cmd="$GOPATH/bin/vrgo --mode=client --port=1000"

for element in ${backups[@]}; do
	IFS=',' read id port <<< "${element}"
	eval "${backup_cmd}  --port=${port} > backup${id}-${port}.log&"
	echo "Running ${backup_cmd}  --port=${port} > backup${id}-${port}.log&"
done

eval ${primary_cmd}
echo "Running ${primary_cmd}"
#eval ${client_cmd}
