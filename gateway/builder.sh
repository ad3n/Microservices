go build -ldflags="-X github.com/devopsfaith/krakend/core.KrakendVersion=KejawenLab" -o api-gateway ./app
go build -buildmode=plugin -o authenticator.so ./plugins/authenticator
