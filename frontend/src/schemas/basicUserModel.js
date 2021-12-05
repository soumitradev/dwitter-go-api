const Joi = require('joi');

const basicUserSchema = Joi.object({
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

  followerCount: Joi.number()
    .integer()
    .required(),

  followingCount: Joi.number()
    .integer()
    .required(),

  createdAt: Joi.date()
    .required(),
}).required();

class basicUserModel {
  constructor(username, name, email, bio, pfpURL, followerCount, followingCount, createdAt) {
    this.username = username;
    this.name = name;
    this.email = email;
    this.bio = bio;
    this.pfpURL = pfpURL;
    this.followerCount = followerCount;
    this.followingCount = followingCount;
    this.createdAt = createdAt;

    var res = basicUserSchema.validate(this);
    if (res.error) {
      throw Error(`basicUserModel schema validation failed: ${res.error}`);
    }
  }
}

export default { basicUserModel, basicUserSchema }