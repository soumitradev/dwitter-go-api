const Joi = require('joi');

import dweetSchema from './dweetModel';
import redweetSchema from './redweetModel';

const feedObjectSchema = Joi.alternatives().try(dweetSchema, redweetSchema);

class feedObjectModel {
  constructor() {
    var res = feedObjectSchema.validate(this);
    if (res.error) {
      throw Error(`feedObjectModel schema validation failed: ${res.error}`);
    }
  }
}

export default { feedObjectModel, feedObjectSchema }