'use strict';

/**
 * Read the documentation (https://strapi.io/documentation/developer-docs/latest/development/backend-customization.html#core-controllers)
 * to customize this controller
 */

 const { sanitizeEntity } = require('strapi-utils');

 module.exports = {
   async find(ctx) {
     let entities;
     if (ctx.query._q) {
       entities = await strapi.services.synapse.search(ctx.query);
     } else {
       entities = await strapi.services.synapse.find(ctx.query);
     }

     entities = entities.map(entity => {
      // const postNeuronPart = entity.postNeuron ? `-${entity.postNeuron?.uid}` : '';
      // const synapseSection = entity.section ? `, section ${entity.section}` : '';
      // const neuronsPostPart = entity.neuronPost.length > 0 ? entity.neuronPost.reduce((r, n, i) => {
      //   let s = (i != 0 ? ', ' : '');
      //   let e = (i == entity.neuronPost.length - 1 ? ')' : '');
      //   return `${r}${s}${n.uid}${e}`
      // }, ' (') : '';
      // return ({
      //    ...entity,
      //    name: entity.position === 'pre' ? `<b>pre-${entity.neuronPre?.uid}</b>-${entity.type}-post${postNeuronPart}${neuronsPostPart}${synapseSection}` : `pre-${entity.neuronPre?.uid}-${entity.type}-<b>post${postNeuronPart}</b>${neuronsPostPart}${synapseSection}`,
      //  });
      return ({
        ...entity,
        name: entity.uid,
      });
      });

     return entities.map(entity => {
       return sanitizeEntity(entity, {
         model: strapi.models.synapse,
       });
     });
   },

   async count(ctx) {
    if (ctx.query.terms != null && ctx.query.terms != "") {
      const terms = ctx.query.terms
        ? ctx.query.terms.split(",").map((t) => t.toLowerCase())
        : [];
      const start = parseInt(ctx.query._start || "0");
      const limit = parseInt(ctx.query._limit || "30");
      const timepoint = ctx.query.timepoint;
      ctx.send(
        await strapi.services.synapse.customSearchCount(timepoint, terms)
      );
    } else {
      if (ctx.query._q) {
        return strapi.services.synapse.countSearch(ctx.query);
      }
      return strapi.services.synapse.count(ctx.query);
    }
   }
 };
