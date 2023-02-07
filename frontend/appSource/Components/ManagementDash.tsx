import * as React from 'react';
import Stack from "@mui/material/Stack";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import TextField from "@mui/material/TextField";
import CardActions from "@mui/material/CardActions";
import Button from "@mui/material/Button";
import {CardContent} from "@mui/material";
import Container from "@mui/material/Container";
import Grid from "@mui/material/Grid";
import ManagementTableConnections from "./ManagementTableConnections";

export default class ManagementDash extends React.Component<any, any> {
    render(): any {
        return (
            <Grid container wrap={'nowrap'} spacing={1} sx={{height: "95%", p: 0, m: 0, width: "100%"}}>
                <Grid item><ManagementTableConnections/></Grid>
                <Grid item>
                    <Container>
                        <Card>
                            <CardContent>
                                <Stack spacing={2}>
                                    <Typography>Добавление нового элемента</Typography>
                                    <TextField id="outlined-basic" label="IP" variant="outlined"/>
                                    <TextField id="outlined-basic" label="Порт" variant="outlined"/>
                                    <TextField id="outlined-basic" label="Имя соединения" variant="outlined"/>
                                    <TextField id="outlined-basic" label="Ключ" variant="outlined"/>
                                </Stack>
                            </CardContent>
                            <CardActions>
                                <Button size="small">Подтвердить</Button>
                            </CardActions>
                        </Card>
                    </Container>
                </Grid>
            </Grid>
        )
    }
}