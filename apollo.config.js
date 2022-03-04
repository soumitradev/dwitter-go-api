module.exports = {
  client: {
    service: {
      name: 'dwitter',
      url: 'http://localhost:5000/api/graphql',
    },
    includes: [
      'frontend/src/**/*.vue',
      'frontend/src/**/*.js',
    ],
  },
}
