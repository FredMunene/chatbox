package handlers

import (
	"forum/backend/models"
)

// Mock repository functions
func mockGetComments(db interface{}, postID int) ([]models.Post, error) {
	return []models.Post{{ID: postID, Body: "Test comment"}}, nil
}

func mockGetReactions(db interface{}, postID int, reactionType string) ([]models.Reaction, error) {
	return []models.Reaction{{ID: postID, Reaction: reactionType}}, nil
}

func mockGetCategories(db interface{}, postID int) ([]models.Category, error) {
	return []models.Category{{ID: postID, CategoryName: "Test Category"}}, nil
}

func mockGetUserByEmail(email string) (models.User, error) {
	return models.User{Username: "TestUser", Email: email}, nil
}

// func TestPostDetails(t *testing.T) {
// 	// Replace repository functions with mocks
// 	repositories.GetComments = mockGetComments
// 	repositories.GetReactions = mockGetReactions
// 	repositories.GetCategories = mockGetCategories
// 	repositories.GetUserByEmail = mockGetUserByEmail

// 	// // Call PostDetails function
// 	// PostDetails(rr, req, posts, false, mockRepo)

// 	// // Check the response status code
// 	// if status := rr.Code; status != http.StatusOK {
// 	// 	t.Errorf("Expected status code %d but got %d", http.StatusOK, status)
// 	// }
// 	// Create a fake HTTP request
// 	req, err := http.NewRequest("GET", "/post-details", nil)
// 	if err != nil {
// 		t.Fatalf("Failed to create request: %v", err)
// 	}

// 	// Create a response recorder
// 	rr := httptest.NewRecorder()

// 	// Mock posts data
// 	posts := []models.Post{
// 		{ID: 1, PostTitle: "Test Post", Body: "This is a test post"},
// 	}

// 	// Call PostDetails function
// 	handlers.PostDetails(rr, req, posts, false)

// 	// Check the response status code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("Expected status code %d but got %d", http.StatusOK, status)
// 	}
// }
