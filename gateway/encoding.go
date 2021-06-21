package krakend

import (
	rss "github.com/devopsfaith/krakend-rss"
	xml "github.com/devopsfaith/krakend-xml"
	ginxml "github.com/devopsfaith/krakend-xml/gin"
	"github.com/devopsfaith/krakend/router/gin"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()

	gin.RegisterRender(xml.Name, ginxml.Render)
}
