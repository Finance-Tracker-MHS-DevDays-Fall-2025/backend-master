package docs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        const spec = %s;
        window.onload = function() {
            SwaggerUIBundle({
                spec: spec,
                dom_id: '#swagger-ui',
            })
        }
    </script>
</body>
</html>`
)

func NewSwaggerHandler(swaggerJSON []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(
				fmt.Sprintf(html, swaggerJSON),
			),
		)
	}
}
