const Joi = require('joi');

const basicDweetSchema = Joi.object({
  dweetBody: Joi.string()
    .max(240)
    .required(),

  id: Joi.string()
    .alphanum()
    .length(10)
    .required(),

  // author: Joi.ref('password'),

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
});

export default basicDweetSchema