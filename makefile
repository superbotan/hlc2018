cdgen:
	easyjson -disallow_unknown_fields -all rawdata.go
build: cdgen
	CGO_ENABLED=0 go build -o bin/app -a -installsuffix cgo ./ 
	docker build -t hlcapi .
run:
	docker run -it --rm -p 8080:80 -m 640m -v /home/scherbina/Documents/highloadcup.ru/go_v0/data/220119/data:/tmp/data hlcapi	
br: build run
runhard:
	docker run -it --rm -p 8080:80 -m 2048m -v /home/scherbina/Documents/highloadcup.ru/go_v0/data/hard/data:/tmp/data hlcapi
brh: build runhard

fight:
	echo "Phase I"
	#highloadcup_tester -addr http://localhost:8080 -hlcupdocs data/291218 -utf8 -phase 1 -concurrent 4 -filter "/accounts/group/" > tank.1.log
	#highloadcup_tester -addr http://localhost:8080 -hlcupdocs ../go_v0/data/291218 -utf8 -phase 1 -concurrent 1 -filter "/accounts/filter/" > tank.1.log
	highloadcup_tester -addr http://localhost:8080 -hlcupdocs ../go_v0/data/220119 -utf8 -phase 1 -concurrent 1 > tank.1.log
	echo "Phase II"
	sleep 5
	highloadcup_tester -addr http://localhost:8080 -hlcupdocs ../GO_v0/data/220119 -utf8 -concurrent 1 -phase 2 > tank.2.log
	echo "Phase III"
	sleep 5
	highloadcup_tester -addr http://localhost:8080 -hlcupdocs ../GO_v0/data/220119 -utf8 -phase 3 -concurrent 4 -filter "/accounts/filter/" > tank.3.log
