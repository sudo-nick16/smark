package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sudo-nick16/smark/galactus/repository"
	"github.com/sudo-nick16/smark/galactus/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SyncBookmarks(userRepo *repository.UserRepo, bookmarksRepo *repository.BookmarksRepo) fiber.Handler {
	failedSyncError := fiber.NewError(fiber.StatusBadRequest, "failed to sync")

	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		authCtx := c.Locals("AuthContext").(*types.AuthTokenClaims)
		if authCtx == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		userId := authCtx.UserId
		var requestBody types.SyncEventRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			return err
		}
		log.Printf("Request body: %v", requestBody)

		log.Printf("Events: %v", requestBody.Events)

		events := requestBody.Events

		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			for _, v := range events {
				switch v.Type {
				case types.CreateListEvent:
					{
						var bl types.BookmarkList
						bl.UserId = userId
						bl.Public = false
						err := mapstructure.Decode(v.Data, &bl)
						log.Printf("Create list event data: %v", bl)
						if err != nil {
							return nil, failedSyncError
						}
						bl.UserId = userId
						list, _ := bookmarksRepo.GetBookmarkListByTitle(bl.Title, userId)
						if list == nil {
							_, err = bookmarksRepo.CreateBookmarkList(&bl)
							if err != nil {
								return nil, failedSyncError
							}
						}
						// continue
						log.Printf("list already exists: %v", bl.Title)
						// return nil, fiber.NewError(fiber.StatusBadRequest, "list already exists")
						break
					}
				case types.UpdateListTitleEvent:
					{
						var data struct {
							OldTitle string `json:"oldTitle"`
							Title    string `json:"title"`
						}
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.UpdateBookmarkListTitle(data.Title, data.OldTitle, userId)
						if err != nil {
							return nil, failedSyncError
						}
						break
					}
				case types.ChangeListVisibilityEvent:
					{
						var data types.BookmarkList
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.ChangeListVisibility(data.Title, userId)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.DeleteListEvent:
					{
						var data types.BookmarkList
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.DeleteBookmarkList(data.Title, userId)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.CreateBookmarkEvent:
					{
						var data types.Bookmark
						data.UserId = userId
						mapstructure.Decode(v.Data, &data)
						_, err := bookmarksRepo.CreateBookmark(&data)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.UpdateBookmarkTitleEvent:
					{
						var data struct {
							OldTitle string `json:"oldTitle"`
							Title    string `json:"title"`
						}
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.UpdateBookmarkTitle(data.Title, data.OldTitle, userId)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.UpdateBookmarkEvent:
					{
						var data struct {
							OldTitle  string `json:"oldTitle"`
							Title     string `json:"title"`
							ListTitle string `json:"listTitle"`
							Url       string `json:"url"`
						}
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.UpdateBookmark(data.Url, data.Title, data.OldTitle, userId, data.ListTitle)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.UpdateBookmarkUrlEvent:
					{
						var data struct {
							Url   string `json:"url"`
							Title string `json:"title"`
						}
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.UpdateBookmarkUrl(data.Url, data.Title, userId)
						if err != nil {
							return nil, err
						}
						break
					}
				case types.DeleteBookmarkEvent:
					{
						var data struct {
							Title     string
							ListTitle string
							UserId    string
						}
						mapstructure.Decode(v.Data, &data)
						err := bookmarksRepo.DeleteBookmark(data.ListTitle, data.Title, userId)
						if err != nil {
							return nil, err
						}
						break
					}
				}
			}
			return nil, nil
		}

		session, err := bookmarksRepo.Client.StartSession()
		if err != nil {
			return err
		}

		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, callback)
		if err != nil {
			return failedSyncError
		}

		return c.JSON(fiber.Map{
			"message": "synced successfully",
		})
	}
}

func GetBookmarks(bookmarksRepo *repository.BookmarksRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := c.Locals("AuthContext").(*types.AuthTokenClaims)
		if authCtx == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		userId := authCtx.UserId
		bookmarks, err := bookmarksRepo.GetBookmarks(userId)
		log.Printf("bookmarks: %v", bookmarks)
		if err != nil {
			return fiber.NewError(fiber.StatusOK, "no bookmarks found")
		}
		return c.JSON(
			fiber.Map{
				"bookmarks": bookmarks,
			},
		)
	}
}

func GetPublicList(bookmarksRepo *repository.BookmarksRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("userId")
		bookmarkListId := c.Params("bookmarkListId")
		if userId == "" || bookmarkListId == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid ids.")
		}
		uid, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user id.")
		}
		id, err := primitive.ObjectIDFromHex(bookmarkListId)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid bookmark list id.")
		}
		bookmarkList, err := bookmarksRepo.GetBookmarkList(id, uid)
		log.Printf("bookmark list: %v", bookmarkList)
		if err != nil {
			return fiber.NewError(fiber.StatusOK, "list not found.")
		}
        if !bookmarkList.Public {
			return fiber.NewError(fiber.StatusOK, "the list is private.")
        }
		return c.JSON(
			fiber.Map{
				"bookmarkList": bookmarkList,
			},
		)
	}
}


func GetShareLink(bookmarksRepo *repository.BookmarksRepo, config *types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := c.Locals("AuthContext").(*types.AuthTokenClaims)
		if authCtx == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		userId := authCtx.UserId
        title := c.Params("title")
        if title == "" {
            return fiber.NewError(fiber.StatusBadRequest, "invalid title")
        }
		bookmarkList, err := bookmarksRepo.GetBookmarkListByTitle(title, userId)
		log.Printf("bookmarks: %v", bookmarkList)
		if err != nil {
			return fiber.NewError(fiber.StatusOK, "no bookmarks found")
		}
        if !bookmarkList.Public {
			return fiber.NewError(fiber.StatusOK, "the list is not public")
        }
        link := fmt.Sprintf("%s/%s/%s", config.ServerUrl, userId, bookmarkList.Id.Hex())
		return c.JSON(
			fiber.Map{
				"shareLink": link,
			},
		)
	}
}

