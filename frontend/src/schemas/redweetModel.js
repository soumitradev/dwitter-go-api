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

export default redweetSchema