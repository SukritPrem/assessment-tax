gobuild:
	/usr/local/go/bin/go build .
	make gorun
gorun:
	/usr/local/go/bin/go run .
test:
	/usr/local/go/bin/go test -v ./Test
refreshbranch:
	git fetch origin           # update all tracking branches, including Branch1
	git rebase origin/main  # rebase on latest main

#Step 1 run command source export.sh
#Step 2 run command make run
# 	source export.sh

run:
	docker compose up -d
	sleep 2 #wait for the container to start
	make gobuild
