import axios from 'axios';
import JSZip from 'jszip';

const download = (fileName, blob) => {
  const url = window.URL.createObjectURL(blob);
  // eslint-disable-next-line no-undef
  const a = document.createElement('a');
  a.style.display = 'none';
  a.href = url;
  a.download = fileName;
  // eslint-disable-next-line no-undef
  document.body.appendChild(a);
  a.click();
  setTimeout(() => {
    // eslint-disable-next-line no-undef
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  }, 100);
  return true;
};

export default {
  zipFiles: [],
  getZipFile(url) {
    const zipFileEl = this.zipFiles.filter((zf) => zf.url === url);
    if (zipFileEl.length !== 0) return zipFileEl[0].zipFile;
    const newZipFile = axios.get(url, {
      responseType: 'arraybuffer',
    }).then((resp) => Buffer.from(resp.data, 'binary'))
      .then((zipContent) => {
        console.log({ zipContent });
        return JSZip().loadAsync(zipContent);
      });

    console.log({ newZipFile });
    this.zipFiles.push({
      url,
      zipFile: newZipFile,
    });

    return newZipFile;
  },
  createZipFile(zipname, files) {
    const zip = JSZip();
    files.forEach((file) => {
      zip.file(`${file.name}`, file.content);
    });
    zip.generateAsync({ type: 'blob' })
      .then((blob) => {
        download(zipname, blob);
      }, (err) => {
        // eslint-disable-next-line no-console
        console.log(err);
      });
    return true;
  },
  async getBase64(zipFile, fileName) {
    console.log({ zipFile });
    console.log({ fileName });
    const zipFileContent = await this.getZipFile(zipFile);
    console.log({ zipFileContent });
    return zipFileContent
      .file(fileName)
      .async('string')
      .then((data) => `data:model/obj;base64,${btoa(data)}`);
  },
};
