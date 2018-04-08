#!/bin/bash

usage() {
	cat <<EOF
Usage: $(basename $0) <command>

Wrappers around core binaries:
    create                Creates DB and Table with users.
	test                  Tests if configuration to DB is correct.
EOF
	exit 1
}


CMD="$1"
MYSQL_ROOT_PASSWORD=$(kubectl get secret --namespace default mysql-mysql -o jsonpath="{.data.mysql-root-password}" | base64 --decode; echo)
MINIKUBE_IP=$(minikube ip)
MYSQL_PORT=$(kubectl get svc -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{"\n"}}{{end}}{{end}}{{end}}')
		
echo "DB connection string: root:$MYSQL_ROOT_PASSWORD@$MINIKUBE_IP:$MYSQL_PORT/users"

shift
case "$CMD" in
	create)
        mysqlsh root:$MYSQL_ROOT_PASSWORD@$MINIKUBE_IP:$MYSQL_PORT --sql -f create.sql
        mysqlsh root:$MYSQL_ROOT_PASSWORD@$MINIKUBE_IP:$MYSQL_PORT/users --sql -f tables.sql
        mysqlsh root:$MYSQL_ROOT_PASSWORD@$MINIKUBE_IP:$MYSQL_PORT/users --sql -f test-data.sql
	;;
	test)
		mysqlsh root:$MYSQL_ROOT_PASSWORD@$MINIKUBE_IP:$MYSQL_PORT/users --sql -e "SELECT * FROM users;"
	;;
	*)
		usage
	;;
esac



