import basicUserModel from './basicUserModel';

const Joi = require('joi');

const basicDweetSchema = Joi.object({
  dweetBody: Joi.string()
    .max(240)
    .required(),

  id: Joi.string()
    .alphanum()
    .length(10)
    .required(),

  author: basicUserModel,

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

  isReply: Joi.boolean()
    .required(),

  originalReplyID: Joi.string()
    .alphanum()
    .length(10),

  replyCount: Joi.number()
    .integer()
    .required(),

  redweetCount: Joi.number()
    .integer()
    .required(),

  media: Joi.array()
    .required()
    .items(Joi.string()),
}).required();

class basicDweetModel {
  constructor(dweetBody, id, author, authorID, postedAt, lastUpdatedAt, likeCount, isReply, originalReplyID, replyCount, redweetCount, media) {
    this.dweetBody = dweetBody;
    this.id = id;
    this.author = author;
    this.authorID = authorID;
    this.postedAt = postedAt;
    this.lastUpdatedAt = lastUpdatedAt;
    this.likeCount = likeCount;
    this.isReply = isReply;
    this.originalReplyID = originalReplyID;
    this.replyCount = replyCount;
    this.redweetCount = redweetCount;
    this.media = media;

    var res = basicDweetSchema.validate(this);
    if (res.error) {
      throw Error(`basicDweetModel schema validation failed: ${res.error}`);
    }
  }
}

export default { basicDweetModel, basicDweetSchema }