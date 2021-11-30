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
});

export default basicUserSchema