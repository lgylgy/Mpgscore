GLIDE = glide

update-deps:
	$(GLIDE) update

install-deps:
	$(GLIDE) install

buildw:
	GOOS=windows GOARCH=amd64 go build ...

buildl:
	GOOS=linux GOARCH=amd64 go build ...