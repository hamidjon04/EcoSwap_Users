// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Hamidjon",
            "email": "nuriddinovhamidjon2@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "user uchun login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User Register",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Xato",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Ushbu endpoint foydalanuvchini tizimdan chiqarish uchun ishlatiladi.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Foydalanuvchini tizimdan chiqarish",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003caccess_token\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Logout successful",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid token or missing Authorization header",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to blacklist access token or cancel token",
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
        "/auth/register": {
            "post": {
                "description": "Register a new user with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User Register",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.UserRegister"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Muvaffaqiyatli ro'yxatdan o'tdingiz",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Xato",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/auth/resetPass": {
            "post": {
                "description": "Reset a user's password with the provided email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "User Email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.Email"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password reset successful",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/updatePass": {
            "post": {
                "description": "Ushbu endpoint foydalanuvchi parolini yangilash uchun ishlatiladi.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Foydalanuvchi parolini yangilash",
                "parameters": [
                    {
                        "description": "Reset Password Payload",
                        "name": "reset_password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ResetPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful Response",
                        "schema": {
                            "$ref": "#/definitions/users.Status"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/users.Status"
                        }
                    }
                }
            }
        },
        "/auth/updateToken": {
            "put": {
                "description": "Ushbu endpoint foydalanuvchi tokenini yangilash uchun ishlatiladi.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Foydalanuvchi tokenini yangilash",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003caccess_token\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access_token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid token or missing Authorization header",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to update token",
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
        "model.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "model.ResetPassword": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                },
                "oldPassword": {
                    "type": "string"
                }
            }
        },
        "users.Email": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "users.Status": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "users.UserLogin": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "users.UserRegister": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Auth Service",
	Description:      "This is the Auth service of EcoSwap project",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
