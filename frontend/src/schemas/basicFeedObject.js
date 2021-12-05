const Joi = require('joi');

import basicDweetSchema from './basicDweetModel';
import redweetSchema from './redweetModel';

const basicFeedObjectSchema = Joi.alternatives().try(basicDweetSchema, redweetSchema);

class basicFeedObjectModel {
  constructor() {
    var res = feedObjectSchema.validate(this);
    if (res.error) {
      throw Error(`feedObjectModel schema validation failed: ${res.error}`);
    }
  }
}

export default { basicFeedObjectModel, basicFeedObjectSchema }