// Импортируем модуль для работ с путями
const path = require('path');
// У самого Webpack уже есть встроенные плагины, их неплохо тоже импортировать
const webpack = require('webpack');
module.exports = {
    watch: true,
    mode: 'development',
    // Указываем входную точку
    entry: './index.ts',
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
      },
    module: {
        rules: [
            { 
              test: /\.less$/i,
              use: [
                // compiles Less to CSS
                "style-loader",
                "css-loader",
                "less-loader",
              ], },
            { test: /\.ts?$/, use: 'ts-loader' }
        ]
      },
      output: {
        // Тут мы указываем полный путь к директории, где будет храниться конечный файл
        path: path.resolve(__dirname, './public/dist'),
        filename: 'bundle.js'
      }
}