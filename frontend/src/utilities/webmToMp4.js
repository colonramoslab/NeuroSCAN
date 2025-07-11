import { backendURL } from './constants';

// Example POST method implementation:
async function postData(url = '', data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
    method: 'POST',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'video/webm',
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
