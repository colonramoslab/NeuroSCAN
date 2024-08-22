'use strict';

const uuid = require('uuid');
const fs = require('fs');
const webmToMp4 = require('webm-to-mp4');

module.exports = {
  webm2avi: async ctx => {
    const { type, data } = ctx.request.body.webmData;

    const id = uuid.v4();
    const inputFile = `/tmp/${id}.webm`;
    const outputFile = `/tmp/${id}.mp4`;
    console.log(`Writing webm to: ${inputFile}`);
    const buf = Buffer.from(data);
    fs.writeFileSync(inputFile, buf);

    await fs.writeFile(outputFile, Buffer.from(webmToMp4(await fs.readFil(inputFile))));

    const aviData = fs.readFileSync(outputFile);

    fs.unlinkSync(inputFile);
    fs.unlinkSync(outputFile);

    return {"result": aviData};
  },
};
