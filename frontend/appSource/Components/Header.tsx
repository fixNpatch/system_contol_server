import * as React from "react";
import {AppBar, Button, IconButton, Toolbar, Typography} from "@mui/material";

export default class Header extends React.Component{
    constructor(props:any) {
        super(props);
        this.state = {open:false};
    }

    render(){
        return (
            <AppBar position="static">
                <Toolbar>
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="menu"
                        sx={{ mr: 2 }}
                    >
                    </IconButton>
                    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                        System Monitor
                    </Typography>
                    <Button color="inherit">Logout</Button>
                </Toolbar>
            </AppBar>
        );
    }
}