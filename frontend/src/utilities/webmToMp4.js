import { backendURL } from './constants';

// Example POST method implementation:
async function postData(url = '', data = {}) {
  // Default options are marked with *
  // set a content-length header to the size of the data
  if (!(data instanceof Blob)) {
    throw new Error('Data must be a Blob');
  }
  const response = await fetch(url, {
    method: 'POST',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'video/webm',
      'Content-Length': data.size.toString(), // set content-length header
    },
    referrerPolicy: 'no-referrer',
    body: data,
  });

  if (!response.ok) {
    throw new Error(`Server error: ${response.status}`);
  }

  return response.json(); // parses JSON response into native JavaScript objects
}

export default async (webmData) => {
  const result = await postData(`${backendURL}videos/webmtomp4`, webmData);
  return result;
};
