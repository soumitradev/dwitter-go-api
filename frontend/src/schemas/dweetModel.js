const Joi = require('joi');

import basicDweetSchema from './basicDweetModel';
import basicUserSchema from './basicUserModel';

const dweetSchema = Joi.object({
  dweetBody: Joi.string()
    .max(240)
    .required(),

  id: Joi.string()
    .alphanum()
    .length(10)
    .required(),

  author: basicUserSchema,

  authorID: Joi.string()
    .alphanum()
    .max(20)
    .required(),

  postedAt: Joi.date()
    .required(),

  lastUpdatedAt: Joi.date()
    .required(),

  likeCount: Joi.number()
    .integer()
    .required(),

  likeUsers: Joi.array()
    .required()
    .items(basicUserSchema),

  isReply: Joi.boolean()
    .required(),

  originalReplyID: Joi.string()
    .alphanum()
    .length(10),

  replyTo: basicDweetSchema,

  replyCount: Joi.number()
    .integer()
    .required(),

  replyDweets: Joi.array()
    .required()
    .items(basicDweetSchema),

  redweetCount: Joi.number()
    .integer()
    .required(),

  redweetDweets: Joi.array()
    .required()
    .items(basicDweetSchema),

  media: Joi.array()
    .required()
    .items(Joi.string()),
});

export default dweetSchema