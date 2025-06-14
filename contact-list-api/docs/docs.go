// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/contacts/": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Lista todos os contatos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Contact"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Cria um novo contato",
                "parameters": [
                    {
                        "description": "Contato",
                        "name": "contact",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HTTPError"
                        }
                    }
                }
            }
        },
        "/contacts/email-providers": {
            "get": {
                "description": "Retorna todos os domínios de e-mail utilizados pelos contatos",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Lista provedores de e-mail",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/contacts/search": {
            "get": {
                "description": "Busca contatos com base em parte do nome",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Busca contatos",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Nome para busca parcial",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Contact"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/contacts/summary": {
            "get": {
                "description": "Obtém estatísticas ou dados agregados sobre os contatos",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Resumo dos contatos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/contacts/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Busca um contato por ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do contato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HTTPError"
                        }
                    }
                }
            },
            "put": {
                "description": "Atualiza os dados de um contato existente",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contacts"
                ],
                "summary": "Atualiza um contato por ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do contato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Dados atualizados do contato",
                        "name": "contact",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Deleta um contato existente usando o ID",
                "tags": [
                    "Contacts"
                ],
                "summary": "Remove um contato",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do contato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.HTTPError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Contact": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "joao@email.com"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "João da Silva"
                },
                "phone": {
                    "type": "string",
                    "example": "11999998888"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Contact List API",
	Description:      "API para gerenciamento de lista de contatos",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
