module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  configureWebpack: {
    resolve: {
      fallback: {
        stream: require.resolve("stream-browserify"),
      },
    }
  }
}