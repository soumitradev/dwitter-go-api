// import basicDweetSchema from './basicDweetModel';
// import basicFeedObject from './basicFeedObject';
// import basicUserSchema from './basicUserModel';
// import redweetSchema from './redweetModel';

const userModel = {
  username: '',
  name: '',
  email: '',
  bio: '',
  pfpURL: '',
  dweets: [],
  redweets: [],
  feedObjects: [],
  likedDweets: [],
  followerCount: 0,
  followers: [],
  followingCount: 0,
  following: [],
  createdAt: '',
}

// class userModel {
//   constructor(username, name, email, bio, pfpURL, dweets, redweets, feedObjects, likedDweets, followerCount, dweets, followingCount, following, createdAt) {
//     this.username = username;
//     this.name = name;
//     this.email = email;
//     this.bio = bio;
//     this.pfpURL = pfpURL;
//     this.dweets = dweets;
//     this.redweets = redweets;
//     this.feedObjects = feedObjects;
//     this.likedDweets = likedDweets;
//     this.followerCount = followerCount;
//     this.followers = followers;
//     this.followingCount = followingCount;
//     this.following = following;
//     this.createdAt = createdAt;

//     var res = userSchema.validate(this);
//     if (res.error) {
//       throw Error(`userModel schema validation failed: ${res.error}`);
//     }
//   }
// }

export default { userModel }