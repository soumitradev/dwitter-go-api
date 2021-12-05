const Joi = require('joi');

import basicDweetSchema from './basicDweetModel';
import basicUserSchema from './basicUserModel';

const redweetSchema = Joi.object({
  author: basicUserSchema,

  authorID: Joi.string()
    .alphanum()
    .max(20)
    .required(),

  redweetOf: basicDweetSchema,

  originalRedweetID: Joi.string()
    .alphanum()
    .length(10)
    .required(),

  redweetTime: Joi.date()
    .required(),
});


class redweetModel {
  constructor(author, authorID, redweetOf, originalRedweetID, redweetTime) {
    this.author = author;
    this.authorID = authorID;
    this.redweetOf = redweetOf;
    this.originalRedweetID = originalRedweetID;
    this.redweetTime = redweetTime;

    var res = redweetSchema.validate(this);
    if (res.error) {
      throw Error(`redweetModel schema validation failed: ${res.error}`);
    }
  }
}

export default { redweetModel, redweetSchema }