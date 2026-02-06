# The NeuroSCAN & PomoterDB Frontend

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app), using
the [Redux](https://redux.js.org/) and [Redux Toolkit](https://redux-toolkit.js.org/) template.

We use [craco](https://www.npmjs.com/package/@craco/craco) to extend the CRA setup.

## Getting Started

This project uses Node version `14.21.3`. Node updates a major version roughly every 6 months so use a tool like [NVM](https://github.com/nvm-sh/nvm) to install the specified version. It is stipulated in the `.nvmrc` file.

We use yarn as package manager. Once you have the correct Node version installed, you can follow these steps to get started:

Install yarn

```bash
npm install -g yarn
```

Install dependencies

```bash
yarn
```

We have to overwrite a couple of vendor files, run the below to do do:

````bash
cp ./overwrite/Canvas.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/
cp ./overwrite/ThreeDEngine.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/
cp ./overwrite/OBJLoader.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/
```

### Configure Backend URL

Environment variable: `REACT_APP_BACKEND_URL`

The backend URL is configured via CRA's built-in `.env` file support. The committed `.env.development` file sets the default for local development:

```
REACT_APP_BACKEND_URL=http://localhost:8123/
```

To override locally without modifying tracked files, create a `.env.development.local` file (gitignored).

In production, when the Go backend serves the frontend build, leave `REACT_APP_BACKEND_URL` unset so API calls use same-origin relative URLs.

### Build

```bash
yarn build
```

This should output static files to the `frontend/build/` directory.
````
