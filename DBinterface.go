package main

func APIGetPost(postID string, replies_to_fetch int) (DweetType, error) {
	post, err := GetFullPost(postID, replies_to_fetch)
	npost := FormatAsDweetType(post)
	return npost, err
}
