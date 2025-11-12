package docs

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const(
	html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                spec: %s,
                dom_id: '#swagger-ui',
            })
        }
    </script>
</body>
</html>`
)

func NewSwaggerRouter(swagger *openapi3.T) echo.HandlerFunc {
	return func(c echo.Context) error {
		specJSON, err := swagger.MarshalJSON()
		if err != nil {
			return err
		}

		return c.HTML(
			http.StatusOK, 
			fmt.Sprintf(html, specJSON),
		)
	}
}
