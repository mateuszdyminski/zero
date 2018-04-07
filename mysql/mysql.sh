#!/bin/bash

usage() {
	cat <<EOF
Usage: $(basename $0) <command> <server-type> <version>

Wrappers around core binaries:
    create                Creates DB and Table with users.
EOF
	exit 1
}


CMD="$1"
MYSQL_ROOT_PASSWORD=$(kubectl get secret --namespace default mysql-mysql -o jsonpath="{.data.mysql-root-password}" | base64 --decode; echo)
		
shift
case "$CMD" in
	create)
        mysqlsh root:$MYSQL_ROOT_PASSWORD@192.168.99.100:31594 --sql -f create.sql
        mysqlsh root:$MYSQL_ROOT_PASSWORD@192.168.99.100:31594/users --sql -f tables.sql
        mysqlsh root:$MYSQL_ROOT_PASSWORD@192.168.99.100:31594/users --sql -f test-data.sql
	;;
	*)
		usage
	;;
esac



