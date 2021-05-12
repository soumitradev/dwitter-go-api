package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func runDBTests() {
	test()
}

func test() error {
	// Bunch of Testing BS, Comment and uncomment based on what you want to test

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("pisscock"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	createdUser, err := NewUser("PyRet", string(passwordHash), "Py", "Ret", "pyret@gmail.com", "le bio")
	if err != nil {
		panic(err)
	}

	// madhav, err := NewUser("Madhav", string(passwordHash), "papa", "peli", "peli@gmail.com", "pr0 at js")
	// if err != nil {
	// 	panic(err)
	// }

	createdPost, err := NewDweet("evening", createdUser.Username, []string{})
	if err != nil {
		panic(err)
	}

	_, err = NewLike(createdPost.ID, createdUser.Username)
	if err != nil {
		panic(err)
	}

	// reply, err := NewReply(createdPost.ID, createdUser.Username, "nice dweet", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// newUser, err := NewUser("YourMom", "Your", "Mom", "mom@gmail.com", "Evening")
	// if err != nil {
	// 	panic(err)
	// }

	// redweet, err := NewRedweet(createdPost.ID, createdUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// updated_post, err := UpdateDweet(createdPost.ID, "mmmm froge", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// updated_user, err := UpdateUser(createdUser.Username, createdUser.Username, createdUser.FirstName, createdUser.LastName, createdUser.Email, "papa peli more liek daddy peli 🥺")
	// if err != nil {
	// 	panic(err)
	// }

	// lepost, err := GetPostReplies(createdPost.ID, 1)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewDweet("getting deleted soon ;)", newUser.Username, []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewLike(createdPost.ID, newUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewReply(createdPost.ID, newUser.Username, "this reply finna die toooo", []string{})
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewRedweet(createdPost.ID, newUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = NewFollower(createdUser.Username, newUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = DeleteFollower(createdUser.Username, newUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// deletedUser, err := DeleteUser(newUser.Username)
	// if err != nil {
	// 	panic(err)
	// }

	// deletedPost, err := DeleteDweet(createdPost.ID)
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
