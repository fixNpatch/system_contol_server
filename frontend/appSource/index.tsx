import App from "./App";
import * as React from 'react'
import {createRoot} from 'react-dom/client';


function ready() {
    console.log("Hello, it's frontend compiled by webpack");

    const rootElement = document.getElementById('root');
    // @ts-ignore
    const root = createRoot(rootElement);

    root.render(React.createElement(App));
}

document.addEventListener("DOMContentLoaded", ready);