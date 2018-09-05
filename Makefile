GLIDE = glide

update-deps:
	$(GLIDE) update

install-deps:
	$(GLIDE) install

buildw:
	GOOS=windows GOARCH=amd64 go install ./...

buildl:
	GOOS=linux GOARCH=amd64 go install ./...