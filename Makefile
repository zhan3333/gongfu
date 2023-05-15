remote_dir = /home/zhan/Application/gongfu
remote_ssh = zhan@t
build_api:
	GOOS=linux GOARCH=amd64 go build -o build/gongfu cmd/main.go
build_web:
	cd web && yarn run build:prod
	rm -rf build/web/*
	cp -r web/dist/angular-ngrx-material-starter/* build/web/

upload_api: build_api
	scp build/gongfu ${remote_ssh}:${remote_dir}/
	scp -r build/config ${remote_ssh}:${remote_dir}/
	scp -r build/storage ${remote_ssh}:${remote_dir}/
	scp build/Dockerfile ${remote_ssh}:${remote_dir}/
	ssh ${remote_ssh} "cd ~/Application && docker-compose up -d --no-deps --build gongfu"

upload_web: build_web
	ssh ${remote_ssh} "rm -rf ${remote_dir}/web/*"
	scp -r build/web/* ${remote_ssh}:${remote_dir}/web/

upload: upload_web upload_api

restart:
	scp -r build/config ${remote_ssh}:${remote_dir}/
	ssh ${remote_ssh} "cd ~/Application && docker-compose up -d --no-deps gongfu"

run:
	go run cmd/main.go

.PHONY: web
web:
	cd web && yarn run start