import path from "path";
import webpack, { experiments } from "webpack";
import nodeExternals from 'webpack-node-externals';
// @ts-ignore
import TypescriptDeclarationPlugin from 'typescript-declaration-webpack-plugin';

const config: webpack.Configuration = {
  entry: "./src/index.tsx",
  devtool: false,
  module: {
    rules: [
      {
        test: /\.(ts|js)x?$/,
        exclude: /(node_modules|bower_components|dist)/,
        use: {
          loader: "babel-loader",
          options: {
            babelrc: true,
          },
        },
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "index.js",
    libraryTarget: 'umd',
    libraryExport: 'default',
  },
  plugins: [
    // new BundleAnalyzerPlugin(),
    new webpack.ProvidePlugin({
      React: 'react',
    }),
    new TypescriptDeclarationPlugin({
      out: 'dist/index.d.ts',
    }),
  ],
  externals: [nodeExternals()],
};

export default config;
