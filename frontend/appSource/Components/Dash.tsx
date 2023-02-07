import * as React from "react";
import {Box, Container, Grid, Paper} from "@mui/material";
import TableConnections from "./TableConnections";
import PlotComponent from "./PlotComponent";
import SystemInfoComponent from "./SystemInfoComponent";
import MiniControlComponent from "./MiniControlComponent";
import {height} from "@mui/system";

const header = 12;

export default class Dash extends React.Component {
    constructor(props:any) {
        super(props);
    }

    render() {
        return (
                <Grid container wrap={'nowrap'} spacing={1} sx={{height:"95%",  p: 0, m:0, width: "100%"}}>
                    <Grid item flex={1} sx={{height:"100%", p: 0, m:0}}>
                        {/*<Paper sx={{height:"100%"}}>s</Paper>*/}
                        <TableConnections/>
                    </Grid>

                    <Grid item flex={1} container direction="column" wrap={'wrap'} spacing={1} sx={{height:"100%", p: 0, m:0, width:"100%"}} >
                        {/*<Grid item sx={{height:"100%"}}>*/}
                        {/*    <Paper sx={{height:"100%", p:0, m:0}}>*/}
                        {/*        <PlotComponent/>*/}
                        {/*    </Paper>*/}
                        {/*    /!*<Paper sx={{height:"100%"}}>s</Paper>*!/*/}
                        {/*</Grid>*/}

                        <Grid item sx={{height:"100%", p: 0, m:0, width: "100%"}}>
                            <Paper sx={{height:"100%", p: 0, m:0, width: "100%", maxWidth: "100%", overflowX:"auto"}}>
                                <MiniControlComponent/>
                            </Paper>
                        </Grid>
                    </Grid>
                </Grid>
            )
    }
}