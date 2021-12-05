const Joi = require('joi');

import basicDweetSchema from './basicDweetModel';
import basicFeedObject from './basicFeedObject';
import basicUserSchema from './basicUserModel';
import redweetSchema from './redweetModel';

const userSchema = Joi.object({
  username: Joi.string()
    .alphanum()
    .max(20)
    .required(),

  name: Joi.string()
    .max(80)
    .required(),

  email: Joi.string()
    .max(100)
    .required(),

  bio: Joi.string()
    .max(160)
    .required(),

  pfpURL: Joi.string()
    .required(),

  dweets: Joi.array()
    .required()
    .items(basicDweetSchema),

  redweets: Joi.array()
    .required()
    .items(redweetSchema),

  feedObjects: Joi.array()
    .required()
    .items(basicFeedObject),

  likedDweets: Joi.array()
    .required()
    .items(basicDweetSchema),

  followerCount: Joi.number()
    .integer()
    .required(),

  followers: Joi.array()
    .required()
    .items(basicUserSchema),

  followingCount: Joi.number()
    .integer()
    .required(),

  following: Joi.array()
    .required()
    .items(basicUserSchema),

  createdAt: Joi.date()
    .required(),
});


class userModel {
  constructor(username, name, email, bio, pfpURL, dweets, redweets, feedObjects, likedDweets, followerCount, dweets, followingCount, following, createdAt) {
    this.username = username;
    this.name = name;
    this.email = email;
    this.bio = bio;
    this.pfpURL = pfpURL;
    this.dweets = dweets;
    this.redweets = redweets;
    this.feedObjects = feedObjects;
    this.likedDweets = likedDweets;
    this.followerCount = followerCount;
    this.followers = followers;
    this.followingCount = followingCount;
    this.following = following;
    this.createdAt = createdAt;

    var res = userSchema.validate(this);
    if (res.error) {
      throw Error(`userModel schema validation failed: ${res.error}`);
    }
  }
}

export default { userModel, userSchema }