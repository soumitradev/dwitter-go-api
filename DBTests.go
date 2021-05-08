package main

import (
	"encoding/json"
	"fmt"
)

func runDBTests() {
	test()
}

func test() error {
	// Bunch of Testing BS, Comment and uncomment based on what you want to test

	createdUser, err := NewUser("PyRet", "Py", "Ret", "pyret@gmail.com", "le bio")
	if err != nil {
		panic(err)
	}

	createdPost, err := NewDweet("just setting up my dwttr", createdUser.DbID, []string{})
	if err != nil {
		panic(err)
	}

	// createdLike, err := NewLike(createdPost.DbID, createdUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// reply, err := NewReply(createdPost.DbID, createdUser.DbID, "nice dweet", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// newUser, err := NewUser("YourMom", "Your", "Mom", "mom@gmail.com", "Evening")
	// if err != nil {
	// 	panic(err)
	// }

	// redweet, err := NewRedweet(createdPost.DbID, createdUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// updated_post, err := UpdateDweet(createdPost.DbID, "mmmm froge", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// updated_user, err := UpdateUser(createdUser.DbID, createdUser.Mention, createdUser.FirstName, createdUser.LastName, createdUser.Email, "papa peli more liek daddy peli 🥺")
	// if err != nil {
	// 	panic(err)
	// }

	// lepost, err := GetPostReplies(createdPost.DbID, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewDweet("getting deleted soon ;)", newUser.DbID, []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewLike(createdPost.DbID, newUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewReply(createdPost.DbID, newUser.DbID, "this reply finna die toooo", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewRedweet(createdPost.DbID, newUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewFollower(createdUser.DbID, newUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = DeleteFollower(createdUser.DbID, newUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// deletedUser, err := DeleteUser(newUser.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	// deletedPost, err := DeleteDweet(createdPost.DbID)
	// if err != nil {
	// 	panic(err)
	// }

	result, _ := json.MarshalIndent(createdUser, "", "  ")
	fmt.Printf("Created User: %s\n", result)

	result1, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("Created Post: %s\n", result1)

	// result2, _ := json.MarshalIndent(createdLike, "", "  ")
	// fmt.Printf("Created Like: %s\n", result2)

	// result5, _ := json.MarshalIndent(reply, "", "  ")
	// fmt.Printf("Reply Dweet: %s\n", result5)

	// result7, _ := json.MarshalIndent(redweet, "", "  ")
	// fmt.Printf("Redweet: %s\n", result7)

	// result8, _ := json.MarshalIndent(updated_post, "", "  ")
	// fmt.Printf("Updated Post: %s\n", result8)

	// result9, _ := json.MarshalIndent(updated_user, "", "  ")
	// fmt.Printf("Updated User: %s\n", result9)

	// result10, _ := json.MarshalIndent(lepost, "", "  ")
	// fmt.Printf("ignore: %s\n", result10)

	// result11, _ := json.MarshalIndent(deletedPost, "", "  ")
	// fmt.Printf("ignore: %s\n", result11)

	// result12, _ := json.MarshalIndent(deletedUser, "", "  ")
	// fmt.Printf("ignore: %s\n", result12)

	return nil
}
