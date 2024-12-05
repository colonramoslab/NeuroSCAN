'use strict';

/**
 * Read the documentation (https://strapi.io/documentation/developer-docs/latest/development/backend-customization.html#core-services)
 * to customize this service
 */

/**
 *
 * This is a rethinking of the original search query that intends to be more simple and efficient.
 * It relies entirely on the uid being a source of information and leveraging that to search for
 * synapses that contain certain neurons in the pre or post position.
 * Ultimately it is simplified, meaning if we are looking for a pre neuron, we can just search for
 * a uid that starts with the pre neuron. If we are looking for a post neuron, we just look for a uid
 * that contains the post neuron.
 */
const searchSynapseByTerms = (params, count = false) => {
  const where = params._where || [];
  const { searchTerms } = where.find((t) => "searchTerms" in t);
  const terms = searchTerms
    .toLowerCase()
    .split("|")
    .map((term) => decodeURIComponent(term));
  const { timepoint } = where.find((t) => "timepoint" in t);

  const type = where.filter((t) => "type" in t).map((t) => t.type) || {
    type: [],
  };

  const { neuronPre } = where.find((t) => "neuronPre" in t) || {
    neuronPre: null,
  };

  const { postNeuron } = where.find((t) => "postNeuron" in t) || {
    postNeuron: null,
  };

  // remove the pre and post neuron from the search terms
  const termsFiltered = terms.filter(
    (term) => term !== neuronPre && term !== postNeuron
  );

  let query = `
    SELECT * FROM synapses
    WHERE timepoint = ${timepoint}
  `

  // If we have search terms, see if the uid contains any of them
  if (termsFiltered.length > 0) {
    query += `
      AND (
    `;

    const termOrs = [];

    termsFiltered.forEach((term) => {
      termOrs.push(`( lower(uid) LIKE '%${term.toLowerCase()}%')`);
    });

    query += termOrs.join(" OR ");

    query += `
      )
    `;
  }

  // If we have a type, filter by it
  if (type.length > 0) {
    query += `
      AND "type" = any(array['${type.join("','")}'])
    `
  }

  // If a neuronPre is present, filter by uid that starts with the neuronPre
  if (neuronPre) {
    query += `
      AND (lower(uid) LIKE '${neuronPre.toLowerCase()}%')
    `;
  }

  // If a postNeuron is present, filter by uid that contains the postNeuron
  // This could be a little redundant with the searchTerms, but it's here for now
  if (postNeuron) {
    query += `
      AND (lower(uid) LIKE '%${postNeuron.toLowerCase()}%')
    `;
  }

  return query;
};

// const searchSynapseByTermsOLD = (
//   params,
// ) => {
//   const where = params._where || [];
//   const { searchTerms } = where.find(t => 'searchTerms' in t);
//   const terms = searchTerms.toUpperCase().split('|').map((term) => decodeURIComponent(term));
//   const { timepoint } = where.find(t => 'timepoint' in t);
//   const type = where.filter(t => 'type' in t).map(t => t.type) || {type: []};
//   const { neuronPre } = where.find(t => 'neuronPre' in t) || {neuronPre: null};
//   const { postNeuron } = where.find(t => 'postNeuron' in t) || {postNeuron: null};
//   const termsPre = neuronPre ? [neuronPre.toUpperCase()] : terms;

//   return termsPre.reduce((r, t, i) => {
//     return `${r} ${i != 0 ? 'UNION ': ''}
//   select *
//   from (
//     select s.*, n_pre.uid as neuronPre_uid
//     from synapses as s
//     join neurons as n_pre on n_pre.id = s."neuronPre"
//     join neurons as n_post on n_post.id = s."postNeuron"
//     where s.timepoint = ${timepoint}
//     and upper(n_pre.uid) like '%${t}%'
//     ${type.length > 0 ? `and s.type in ('${type.join("','")}')` : ''}
//     ${(postNeuron && !neuronPre) ? `and upper(n_post.uid) like '%${postNeuron.toUpperCase()}%'` : ''}
//     ${(!postNeuron && neuronPre) ? `and upper(n_pre.uid) like '%${neuronPre.toUpperCase()}%' and s.position = 'pre'` : ''}
//     and ${terms.length - 1} <= (
//       select count(1)
//       from (
//         ${terms.reduce((r2, term, i2) => {
//           let a = (i2 != 0 ? 'union all ' : '');
//           return `${r2} ${a} select 1
//             from synapses__neuron_post snp
//             join neurons n on n.id = snp.neuron_id
//             where upper(n.uid) like '%${term}%'
//             and snp.synapse_id = s.id`
//           }, '')}
//         )
//       )
//   ${!neuronPre ? `
//     union
//       select s.*, n_pre.uid as "neuronPre_uid"
//       from synapses as s
//       join neurons as n_pre on n_pre.id = s."neuronPre"
//       join neurons as n_post on n_post.id = s."postNeuron"
//       where s.timepoint = ${timepoint}
//       ${type.length > 0 ? `and type in ('${type.join("','")}')` : ''}
//       ${postNeuron ? `and upper(n_post.uid) like '%${postNeuron.toUpperCase()}%'` : ''}
//       and ${terms.length} = (
//         select count(1)
//         from (
//           ${terms.reduce((r2, term, i2) => {
//             let a = (i2 != 0 ? 'union all ' : '');
//             return `${r2} ${a} select 1
//               from synapses__neuron_post snp
//               join neurons n on n.id = snp.neuron_id
//               where upper(n.uid) like '%${term}%'
//               and snp.synapse_id = s.id`
//             }, '')}
//           )
//         )
//   ` : ''}) `}, '');
// }

module.exports = {
  async find(params, populate) {
    const where = params._where || [];
    const searchTerms = where.find(t => 'searchTerms' in t);
    if (!searchTerms) {
      return strapi.query('synapse').find(params, populate);
    }

    const limit = params._limit || null;
    const offset = params._start || null;

    const query = `
    select *
    from (
    ${searchSynapseByTerms(params)}
    )
    ${limit !== null ? `limit (${limit})` : ''}
    ${offset !== null ? `offset (${offset})` : ''}
    `;
    const knex = strapi.connections.default;
    const ids = await knex.raw(query);

    const rows = ids.rows || [];

    return await strapi.query('synapse').find({ id_in: rows.map(x => x.id), _sort: 'uid' });
  },

  async count(params, populate) {
    const where = params._where || [];
    const searchTerms = where.find(t => 'searchTerms' in t);
    if (!searchTerms) {
      return strapi.query('synapse').count(params, populate);
    }
    const query = `
    select count(1) as c
    from (
    ${searchSynapseByTerms(params)}
    )
    `
    const knex = strapi.connections.default;
    // log the query for debugging
    const r = await knex.raw(query);

    const rows = r.rows || [];

    return rows[0].c;
  }
};
