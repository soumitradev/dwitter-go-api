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

  redweetUsers: Joi.array()
    .required()
    .items(basicUserSchema),

  media: Joi.array()
    .required()
    .items(Joi.string()),
});

class dweetModel {
  constructor(dweetBody, id, author, authorID, postedAt, lastUpdatedAt, likeCount, likeUsers, isReply, originalReplyID, replyTo, replyCount, replyDweets, redweetCount, redweetUsers, media) {
    this.dweetBody = dweetBody;
    this.id = id;
    this.author = author;
    this.authorID = authorID;
    this.postedAt = postedAt;
    this.lastUpdatedAt = lastUpdatedAt;
    this.likeCount = likeCount;
    this.likeUsers = likeUsers;
    this.isReply = isReply;
    this.originalReplyID = originalReplyID;
    this.replyTo = replyTo;
    this.replyCount = replyCount;
    this.replyDweets = replyDweets;
    this.redweetCount = redweetCount;
    this.redweetUsers = redweetUsers;
    this.media = media;

    var res = dweetSchema.validate(this);
    if (res.error) {
      throw Error(`dweetModel schema validation failed: ${res.error}`);
    }
  }
}

export default { dweetModel, dweetSchema }