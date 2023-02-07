import * as React from 'react';
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import Table from "@mui/material/Table";
import {connect} from "react-redux";
import { w3cwebsocket as WebSocket3 } from "websocket";
import TableSortLabel from "@mui/material/TableSortLabel";
import Box from "@mui/material/Box";
import {visuallyHidden} from "@mui/utils";
import {Button} from "@mui/material";
import * as $ from 'jquery';

const mapStateToProps = (state:any) => {
    return {
        observedId: state
    };
};


interface Data {
    FakeId: number;
    LAddr: string;
    LPort: number;
    RAddr: string;
    RPort: number;
    PID: number;
    ProcName: string;
    ProcOwner: string;
}

function createData(
    FakeId: number,
    LAddr: string,
    LPort: number,
    RAddr: string,
    RPort: number,
    PID: number,
    ProcName: string,
    ProcOwner: string,
): Data {
    return {
        FakeId,
        LAddr,
        LPort,
        RAddr,
        RPort,
        PID,
        ProcName,
        ProcOwner,
    };
}


function descendingComparator<T>(a: T, b: T, orderBy: keyof T) {
    if (b[orderBy] < a[orderBy]) {
        return -1;
    }
    if (b[orderBy] > a[orderBy]) {
        return 1;
    }
    return 0;
}

type Order = 'asc' | 'desc';

function getComparator<Key extends keyof any>(
    order: Order,
    orderBy: Key,
): (
    a: { [key in Key]: number | string },
    b: { [key in Key]: number | string },
) => number {
    return order === 'desc'
        ? (a, b) => descendingComparator(a, b, orderBy)
        : (a, b) => -descendingComparator(a, b, orderBy);
}

interface HeadCell {
    disablePadding: boolean;
    id: keyof Data;
    label: string;
    numeric: boolean;
}

const headCells: readonly HeadCell[] = [
    {
        id: 'LAddr',
        numeric: false,
        disablePadding: false,
        label: 'LAddr',
    },
    {
        id: 'LPort',
        numeric: true,
        disablePadding: false,
        label: 'LPort',
    },
    {
        id: 'RAddr',
        numeric: false,
        disablePadding: false,
        label: 'RAddr',
    },
    {
        id: 'RPort',
        numeric: true,
        disablePadding: false,
        label: 'RPort',
    },
    {
        id: 'PID',
        numeric: true,
        disablePadding: false,
        label: 'Pid',
    },
    {
        id: 'ProcName',
        numeric: false,
        disablePadding: false,
        label: 'Название процесса',
    },
    {
        id: 'ProcOwner',
        numeric: false,
        disablePadding: false,
        label: 'Имя пользователя',
    },
];

interface EnhancedTableProps {
    onRequestSort: (event: React.MouseEvent<unknown>, property: keyof Data) => void;
    order: Order;
    orderBy: string;
    rowsNumber: number
}



function EnhancedTableHead(props: EnhancedTableProps) {

    const {order, orderBy, onRequestSort, rowsNumber} =
        props;
    const createSortHandler =
        (property: keyof Data) => (event: React.MouseEvent<unknown>) => {
            onRequestSort(event, property);
        };

    return (
        <TableHead>
            <TableRow>
                <TableCell>{rowsNumber}</TableCell>
                {headCells.map((headCell) => (
                    <TableCell
                        key={headCell.id}
                        align={headCell.numeric ? 'right' : 'center'}
                        padding={headCell.disablePadding ? 'none' : 'normal'}
                        sortDirection={orderBy === headCell.id ? order : false}
                    >
                        <TableSortLabel
                            active={orderBy === headCell.id}
                            direction={orderBy === headCell.id ? order : 'asc'}
                            onClick={createSortHandler(headCell.id)}
                        >
                            {headCell.label}
                            {orderBy === headCell.id ? (
                                <Box component="span" sx={visuallyHidden}>
                                    {order === 'desc' ? 'sorted descending' : 'sorted ascending'}
                                </Box>
                            ) : null}
                        </TableSortLabel>
                    </TableCell>
                ))}
            </TableRow>
        </TableHead>
    );
}

class ConnectionsListComponent extends React.Component<any, any> {
    private client: any;
    state: {
        sav: any
        order: Order
        orderBy: keyof Data
        data: any
    };

    constructor(props:any) {
        super(props);
        this.state = {
            order: "asc",
            orderBy: "LAddr",
            sav: null,
            data: [],
        };
    }

    componentDidMount(): void {
        this.setState({
            data: [],
        });
    }

    componentWillUnmount = ():void => {
        this.closeConnection();
    };

    componentDidUpdate = (prevProps: Readonly<any>, prevState: Readonly<any>, snapshot?: any): void => {
        if (this.props.observedId === undefined ) { // если state === undefined => снята галочка
            this.closeConnection(); // закрываем существующее соединение
            return // выходим из функции. Она будет перезапущена через один тик.
        } else if (this.state.sav != this.props.observedId) { // иначе, проверяем не было ли переключено соединение
            this.closeConnection(); // если да, то закрываем "старое".
        }
        if (this.client == null) {
            this.createConnection(this.props.observedId);
        }
    };

    closeConnection = ():void => {
        if (this.client != null) {
            this.client.close(1000, "finish update host info");
        }
        this.client = null;
    };

    createConnection = (id:any):void => {
        this.setState({
            sav: this.props.observedId
        });
        this.client = new WebSocket3('ws://localhost/update/network?ip=' + id);
        this.client.onmessage = (message) => {
            let data = JSON.parse(message.data);
            let connections = data.Connections;
            console.log(connections);

            this.state.data = [];

            connections.forEach((elem, idx, array) => {
                this.state.data.push(createData(elem.FakeId, elem.LAddr.join("."), elem.LPort, elem.RAddr.join("."), elem.RPort, elem.Pid, elem.ProcName, elem.ProcOwner));
            });

            this.forceUpdate();
            this.client.send("ok");
        };

        this.client.onclose = function(event) {
            if (event.wasClean) {
                console.warn(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                console.warn('[close] Connection died');
            }
        };

        this.client.onerror = function(error) {
            console.warn(`[error] ${error.message}`);
        };
    };

    handleRequestSort = (
        event: React.MouseEvent<unknown>,
        property: keyof Data,
    ) => {
        const isAsc = this.state.orderBy === property && this.state.order === 'asc';
        this.setState({
            order: isAsc ? 'desc' : 'asc',
            orderBy: property,
        });
    };


    handleKillConnection(connectionId:any):any {
        console.warn("Connection killed " + connectionId);
        $.post("http://localhost/close", {fakeId:connectionId});
    }

    render(){
        return (
            <TableContainer sx={{
                width: '100%',
                bgcolor: 'background.paper',
                position: 'relative',
                overflow: 'auto',
                maxHeight:'100%',
                maxWidth: '46.07vw',
                fontSize: 10,
                '& ul': { padding: 0 },
            }}>
                <Table
                    sx={{ minWidth: 650, "& th": {
                            // color: "rgba(96, 96, 96)",
                            backgroundColor: '#565f6e',
                            color: "white"
                        }}} size="small" aria-label="List of followed hosts" stickyHeader
                >
                    <EnhancedTableHead
                        order={this.state.order}
                        orderBy={this.state.orderBy}
                        onRequestSort={this.handleRequestSort}
                        rowsNumber={this.state.data.length}
                    />
                    <TableBody>
                        {this.state.data.sort(getComparator(this.state.order, this.state.orderBy))
                            .map((row, index) => {
                                return (
                                    <TableRow
                                        hover
                                        role="checkbox"
                                        tabIndex={-1}
                                        key={row.FakeId}
                                    >
                                        <TableCell><Button color="error" onClick={() => {this.handleKillConnection(row.FakeId)}}>KILL</Button></TableCell>
                                        <TableCell
                                            scope="row"
                                            padding="none"
                                            align="center"
                                        >
                                            {row.LAddr}
                                        </TableCell>
                                        <TableCell align="right">{row.LPort}</TableCell>
                                        <TableCell align="center">{row.RAddr}</TableCell>
                                        <TableCell align="right">{row.RPort}</TableCell>
                                        <TableCell align="right">{row.PID}</TableCell>
                                        <TableCell align="left">{row.ProcName}</TableCell>
                                        <TableCell align="left">{row.ProcOwner}</TableCell>
                                    </TableRow>
                                );
                            })}
                    </TableBody>
                </Table>
            </TableContainer>
        )};
}

export default connect(mapStateToProps)(ConnectionsListComponent);