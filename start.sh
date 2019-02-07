#!/bin/bash
nohup redis-server &
uwsgi --http 0.0.0.0:8000 --module mymodule.wsgi