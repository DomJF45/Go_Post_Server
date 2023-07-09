package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"post_server/configs"
	"post_server/models"
	"post_server/responses"
)

type ProjectController struct {
	Path string
}

var (
	projectCollection *mongo.Collection = configs.GetCollection(configs.DB, "projects")
	validate                            = validator.New()
)

func (c *ProjectController) InitRouter(a *fiber.App) {
	a.Post(c.Path+"/", createPost)          // creates a post
	a.Get(c.Path+"/", getAllPosts)          // gets all posts
	a.Get(c.Path+"/:postId", getPosts)      // returns a post
	a.Delete(c.Path+"/:postId", deletePost) // deletes a post
	a.Put(c.Path+"/:postId", editProject)   // edits a project
}

func NewProjectController() *ProjectController {
	return &ProjectController{
		Path: "/project",
	}
}

func createPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var post models.ProjectPageData
	defer cancel()

	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&post); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	newProject := models.ProjectPageData{
		Id:          primitive.NewObjectID(),
		ProjectName: post.ProjectName,
		Posts:       post.Posts,
		Stack:       post.Stack,
	}

	result, err := projectCollection.InsertOne(ctx, newProject)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusCreated).JSON(responses.ProjectResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func getPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postId := c.Params("postId")
	var posts models.ProjectPageData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := projectCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&posts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": posts},
	})
}

func getAllPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var allPosts []models.ProjectPageData
	defer cancel()

	results, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singlePost models.ProjectPageData
		if err = results.Decode(&singlePost); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
		allPosts = append(allPosts, singlePost)
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": allPosts},
	})
}

func deletePost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postId := c.Params("pageId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	result, err := projectCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.ProjectResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func editProject(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	projectId := c.Params("projectId")
	var project models.ProjectPageData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(projectId)

	if err := c.BodyParser(&project); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&project); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ProjectResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	update := bson.M{"projectName": project.ProjectName, "posts": project.Posts, "stack": project.Stack}

	result, err := projectCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	var updatedProject models.ProjectPageData

	if result.MatchedCount == 1 {
		err := projectCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedProject)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.ProjectResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.ProjectResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": updatedProject},
	})
}
