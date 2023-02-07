import * as React from 'react';
import {List, ListItem, ListItemText, ListSubheader} from "@mui/material";
import {fontSize} from "@mui/system";


export default class ProcessListComponent extends React.Component {
    render() {
        return (

            <List sx={{
                width: '100%',
                // maxWidth: 360,
                bgcolor: 'background.paper',
                position: 'relative',
                overflow: 'auto',
                maxHeight: 300,
                '& ul': { padding: 0 },
            }}
            >
                <ListItem key={`proc-0`}>
                    <ListItemText primary={`Proc0`} />
                </ListItem>
                <ListItem key={`proc-1`}>
                    <ListItemText primary={`Proc1`} />
                </ListItem>
                <ListItem key={`proc-2`}>
                    <ListItemText primary={`Proc2`} />
                </ListItem>
                <ListItem key={`proc-3`}>
                    <ListItemText primary={`Proc3`} />
                </ListItem>
                <ListItem key={`proc-4`}>
                    <ListItemText primary={`Proc4`} />
                </ListItem>
                <ListItem key={`proc-5`}>
                    <ListItemText primary={`Proc5`} />
                </ListItem>
                <ListItem key={`proc-6`}>
                    <ListItemText primary={`Proc6`} />
                </ListItem>
                <ListItem key={`proc-7`}>
                    <ListItemText primary={`Proc7`} />
                </ListItem>
                <ListItem key={`proc-8`}>
                    <ListItemText primary={`Proc8`} />
                </ListItem>
                <ListItem key={`proc-9`}>
                    <ListItemText primary={`Proc9`} />
                </ListItem>
            </List>

        );
    }
}