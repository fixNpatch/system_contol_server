import * as React from "react";
import {Box, Button, Container, Grid, Paper} from "@mui/material";
import { w3cwebsocket as WebSocket3 } from "websocket";
import * as $ from 'jquery';


export default class Settings extends React.Component {
    // define props and state
    props:{};
    state:{
        mounted: boolean,
        someVariable: any,
    };
    private client;

    // initial variables
    constructor(props:any) {
        super(props);

        this.state = {
            mounted: false,
            someVariable: "not set",
        };

        console.log("init settings");
    }


    componentDidMount(): void {
        this.state.mounted = true;
        this.client = new WebSocket3('ws://localhost/update/test');

        this.client.onmessage = (message) => {
            console.log(message);
            this.setState({
                someVariable: message.data,
            });
            setTimeout(() => {this.client.send("ok")}, 1000);
        };

        this.client.onclose = function(event) {
            if (event.wasClean) {
                console.warn(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                // e.g. server process killed or network down
                // event.code is usually 1006 in this case
                console.warn('[close] Connection died');
            }
        };

        this.client.onerror = function(error) {
            console.warn(`[error] ${error.message}`);
        };
    }


    componentWillUnmount(): void {
        console.log("unmount");
        this.client.close(1000, "finish update")
    }

    clearAllData():any {
        $.post("http://localhost/clearall")
    }


    render() {
        return (<div>
            {this.state.someVariable}
            <Button variant="contained" color="success" onClick={() => {
                this.client.close(1000, "component unmount");
                // this.client.close()
            }}>Test</Button>
            <Button variant="contained" color="error" onClick={() => {
                this.clearAllData();
            }}>Clear All Data</Button>
        </div>)
    }
}