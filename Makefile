init_tests:
	go get -u -v github.com/mailcourses/technopark-dbms-forum@master
	go build github.com/mailcourses/technopark-dbms-forum

fill: init_tests
	./technopark-dbms-forum fill --url=http://localhost:8080/api --timeout=900

perf: init_tests
	./technopark-dbms-forum perf --url=http://localhost:8080/api --duration=600 --step=60