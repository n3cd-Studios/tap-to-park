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
        "/admin/organization": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get all of the organizations associated with an admin",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/database.Organization"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create an invite to allow new user to join admin's organization",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Invite"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User or Organization not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to create invite",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/organization/data": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get all of the spots data associated with an organization",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/database.Spot"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/info": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the info of the current user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.JWTResponse"
                        }
                    },
                    "400": {
                        "description": "Failed to use token to retrieve user information",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
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
                            "$ref": "#/definitions/routes.JWTResponse"
                        }
                    },
                    "400": {
                        "description": "Failed to log in",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Registers a User in using a username and a password",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.JWTResponse"
                        }
                    },
                    "400": {
                        "description": "Failed to register user",
                        "schema": {
                            "type": "string"
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
        "/spots/create": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a spot at a longitude and latitude",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Spot"
                        }
                    },
                    "400": {
                        "description": "Invalid body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/spots/delete": {
            "delete": {
                "summary": "Delete a spot by it's ID",
                "responses": {
                    "200": {
                        "description": "Successfully deleted spot",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/spots/info": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get the spots near a longitude and latitude",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Spot"
                        }
                    },
                    "404": {
                        "description": "Spot was not found",
                        "schema": {
                            "type": "string"
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
                "summary": "Get the spots near a longitude and latitude, with optional handicap filter",
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
                    },
                    {
                        "type": "boolean",
                        "description": "filter spots by handicap accessibility",
                        "name": "handicap",
                        "in": "query"
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
                        "description": "Could not load the list of spots",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "database.Coordinates": {
            "type": "object",
            "required": [
                "latitude",
                "longitude"
            ],
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "database.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "database.Invite": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "createdBy": {
                    "type": "integer"
                },
                "expiration": {
                    "type": "string"
                },
                "organization": {
                    "type": "integer"
                },
                "usedBy": {
                    "type": "integer"
                }
            }
        },
        "database.Organization": {
            "type": "object",
            "properties": {
                "invites": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Invite"
                    }
                },
                "name": {
                    "type": "string"
                },
                "spots": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Spot"
                    }
                }
            }
        },
        "database.Reservation": {
            "type": "object",
            "properties": {
                "costPerHour": {
                    "type": "integer"
                },
                "end": {
                    "type": "string"
                },
                "guid": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "start": {
                    "type": "string"
                }
            }
        },
        "database.Spot": {
            "type": "object",
            "properties": {
                "coords": {
                    "$ref": "#/definitions/database.Coordinates"
                },
                "guid": {
                    "type": "string"
                },
                "handicap": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "organization": {
                    "type": "integer"
                },
                "reservations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Reservation"
                    }
                }
            }
        },
        "routes.JWTResponse": {
            "type": "object",
            "properties": {
                "token": {
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
