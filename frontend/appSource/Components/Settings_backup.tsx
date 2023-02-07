// import * as React from "react";
// import {Box, Button, Container, Grid, Paper} from "@mui/material";
// import { w3cwebsocket as WebSocket3 } from "websocket";
//
// const header = 12;
// const client = new WebSocket3('ws://localhost/update');
//
//
// export default class Settings extends React.Component {
//     private Variable:any;
//
//     constructor(props:any) {
//         super(props);
//
//         this.Variable = "not started";
//
//         this.update()
//     }
//
//     private update() {
//
//
//         console.log("try to open websocket");
//
//         // client.send("ready");
//
//         client.onopen = () => {
//             console.log('WebSocket Client Connected');
//         };
//
//         client.onmessage = (message) => {
//             this.Variable = message.data
//         };
//
//         client.onerror = () => {
//             this.Variable = "Stopped"
//         }
//     }
//
//
//     render() {
//         return (<div>
//             {this.Variable}
//             <Button variant="contained" color="success" onClick={() => {
//             }}>Test</Button>
//         </div>)
//     }
// }