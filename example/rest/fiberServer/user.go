package fiberserver

import "github.com/gofiber/fiber/v3"

type userServer struct {
}

func NewUserServer() *userServer {
	return &userServer{}
}

func (u *userServer) GetUserInfo(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": fiber.Map{
			"name": "taluos",
			"age":  20,
		},
	})
}
