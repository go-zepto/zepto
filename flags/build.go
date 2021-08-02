package flags

// Flag used to identify the mode of build when compiling: development/production
// This is used to restrict some development features and create a clean build for production.
// When BuildMode="production" the Zepto server will execute the webpack in production mode
var BuildMode = "development"
