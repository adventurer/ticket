#!/bin/sh
ab -v4 -c10 -n1000 -mGET 'http://127.0.0.1:8088/ticket?c=1'&
ab -v4 -c10 -n1000 -mGET 'http://127.0.0.1:8088/ticket?c=2'&
ab -v4 -c10 -n1000 -mGET 'http://127.0.0.1:8088/ticket?c=3'&
ab -v4 -c10 -n1000 -mGET 'http://127.0.0.1:8088/ticket?c=4'&
ab -v4 -c10 -n1000 -mGET 'http://127.0.0.1:8088/ticket?c=5'&