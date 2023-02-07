import * as React from "react";
import {Box, Button, Container, Grid, Paper, Tab, Tabs, Typography} from "@mui/material";
import ProcessListComponent from "./ProcessListComponent";
import SystemInfoComponent from "./SystemInfoComponent";
import ConnectionsListComponent from "./ConnectionsListComponent";
import MiniCommandComponent from "./MiniCommandComponent";
import {TabContext, TabList, TabPanel} from '@mui/lab';


export default function MiniControlComponent(){
    const [value, setValue] = React.useState('1');

    const handleChange = (event: React.SyntheticEvent, newValue: string) => {
        setValue(newValue);
    };

    return (
        <TabContext value={value}>
            <Box sx={{  p: 0, m:0, borderBottom: 1, borderColor: 'divider', height:"5.8%", width: "100%"}}>
                <TabList onChange={handleChange} aria-label="lab API tabs example">
                    <Tab label="Характеристики" value="1" />
                    <Tab label="Сеть" value="2" />
                    {/*<Tab label="Процессы" value="3" />*/}
                    <Tab label="Управление" value="4" />
                </TabList>
            </Box>
            <TabPanel value="1" sx={{p:0, height:"94.2%", width: "100%"}}><SystemInfoComponent/></TabPanel>
            <TabPanel value="2" sx={{p:0, height:"94.2%", width: "100%", maxWidth:"100%"}}><ConnectionsListComponent/></TabPanel>
            {/*<TabPanel value="3" sx={{p:0}}><ProcessListComponent/></TabPanel>*/}
            <TabPanel value="4" sx={{p:0, height:"94.2%", width: "100%"}}><MiniCommandComponent/></TabPanel>
        </TabContext>
        )
}