const path = require('path');

module.exports = {
    entry: {
        App: './frontend/appSource/index.tsx',
        // Auth: './modules/front/Auth/index.tsx'
    },
    devtool: 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/
            }
        ]
    },
    resolve: {
        extensions: [ '.tsx', '.ts', '.js' ]
    },
    output: {
        filename: 'static/src/modules/[name]/bundle.js',
        path: __dirname,
    },
};