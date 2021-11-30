const Joi = require('joi');

import basicDweetSchema from './basicDweetModel';
import basicUserSchema from './basicUserModel';

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

export default userSchema