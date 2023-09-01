remote_dir = /home/zhan/Application/gongfu
remote_ssh = zhan@t

# 编译 api 服务
build_api:
	cd server && GOOS=linux GOARCH=amd64 go build -o ../build/gongfu cmd/main.go

# 编译 web 服务
build_web:
	cd web && pnpm run build:prod
	rm -rf build/web/*
	cp -r web/dist/gongfu/* build/web/

# 更新 api 服务
upload_api: build_api
	scp build/gongfu ${remote_ssh}:${remote_dir}/
	scp -r build/config ${remote_ssh}:${remote_dir}/
	scp build/Dockerfile ${remote_ssh}:${remote_dir}/
	ssh ${remote_ssh} "cd ~/Application && docker-compose up -d --no-deps --build gongfu"

# 更新 web 服务
upload_web: build_web
	ssh ${remote_ssh} "rm -rf ${remote_dir}/web/*"
	scp -r build/web/* ${remote_ssh}:${remote_dir}/web/

# 更新 web & api 服务
upload: upload_web upload_api

# 重启服务
restart:
	scp -r build/config ${remote_ssh}:${remote_dir}/
	ssh ${remote_ssh} "cd ~/Application && docker-compose up -d --no-deps gongfu"

.PHONY: web
web:
	cd web && yarn run start