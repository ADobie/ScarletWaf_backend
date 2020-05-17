.PHONY: build all deploy cleanall


buildall:
	GOOS=linux GOARCH=amd64 go build . -o scarlet_linux_amd64
	GOOS=darwin GOARCH=amd64 go build . -o scarlet_darwin_amd64
	GOOS=windows GOARCH=amd64 go build . -o scarlet_windows_amd64



buildlinux:
	GOOS=linux GOARCH=amd64 go build  -o scarlet_linux_amd64 .

deploy: buildlinux
	@echo "scp ./scarlet_linux_amd64 scarlet@39.96.77.139:/usr/local/src/scarletBackend/scarlet_linux_amd64" | pbcopy
	@echo "[+] Cmd has copied to the pastebin"
	@echo "[+] Now You can ssh waf and run 'make deploy' in the project dir"

cleanall:
	# clean all executable file
	ls -l  | grep -v "total"|  awk '{print $(NF)}' | xargs file | grep -E "executable" | awk -F:  '{print $1}'| xargs -t rm
