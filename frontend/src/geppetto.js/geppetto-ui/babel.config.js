module.exports = {
  "presets": [
    '@babel/preset-env',
    "@babel/preset-react",
    "@babel/preset-typescript"
  ],
  "plugins": [
    "@babel/plugin-transform-regenerator",
    ["@babel/plugin-proposal-class-properties", { loose: false }],
  ]
}
