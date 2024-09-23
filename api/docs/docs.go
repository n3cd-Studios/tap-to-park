// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://n3cd.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Logs a User in using a username and a password",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.LoginInput"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/database.Error"
                        }
                    }
                }
            }
        },
        "/reservations": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Creates a reservation using a spotID and specified time",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ReservationInput"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/database.Error"
                        }
                    }
                }
            }
        },
        "/reservations/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get a reservation by an ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the reservation",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Spot"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/database.Error"
                        }
                    }
                }
            }
        },
        "/spots/near": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get the spots near a longitude and latitude",
                "parameters": [
                    {
                        "type": "number",
                        "description": "latitude to search by",
                        "name": "lat",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "longitude to search by",
                        "name": "lng",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Spot"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/database.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "database.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "database.Point": {
            "type": "object",
            "properties": {
                "x": {
                    "type": "number"
                }
            }
        },
        "database.Spot": {
            "type": "object",
            "properties": {
                "coords": {
                    "$ref": "#/definitions/database.Point"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "handicap": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "organizationID": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "gorm.DeletedAt": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        },
        "routes.LoginInput": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "routes.ReservationInput": {
            "type": "object",
            "properties": {
                "spotID": {
                    "type": "integer"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Tap-To-Park API",
	Description:      "This is the API for interacting with internal Tap-To-Park services",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
