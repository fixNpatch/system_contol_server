import * as React from 'react';
import Box from '@mui/material/Box';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TableSortLabel from '@mui/material/TableSortLabel';
import Paper from '@mui/material/Paper';
import Checkbox from '@mui/material/Checkbox';
import { visuallyHidden } from '@mui/utils';
import { w3cwebsocket as WebSocket3 } from "websocket";
import { connect } from "react-redux";
import {setObservedConnection} from "../redux/actions";
import {getComparator} from "../comparators";


function createData(
    id: number,
    name: string,
    IP: string,
    status: number,
): Data {
    return {
        id,
        name,
        IP,
        status,
    };
}

const headCells: readonly HeadCell[] = [
    {
        id: 'IP',
        numeric: false,
        disablePadding: true,
        label: 'IP',
    },
    {
        id: 'name',
        numeric: true,
        disablePadding: false,
        label: 'Название',
    },
    {
        id: 'status',
        numeric: true,
        disablePadding: false,
        label: 'Статус соединения',
    },
];

interface EnhancedTableProps {
    onRequestSort: (event: React.MouseEvent<unknown>, property: keyof Data) => void;
    order: Order;
    orderBy: string;
}



function EnhancedTableHead(props: EnhancedTableProps) {
    const {order, orderBy, onRequestSort } =
        props;
    const createSortHandler =
        (property: keyof Data) => (event: React.MouseEvent<unknown>) => {
            onRequestSort(event, property);
        };

    return (
        <TableHead>
            <TableRow>
                <TableCell/>
                {headCells.map((headCell) => (
                    <TableCell
                        key={headCell.id}
                        align={headCell.numeric ? 'right' : 'left'}
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

const mapDispatchToProps = (dispatch:any) => {
    return {
        setId: (id:any) => {
            dispatch(setObservedConnection(id))
        }
    }
};

class TableConnections extends React.Component<any, any>{
    private client: any;
    // private data:any;
    state: {
        order: Order
        orderBy: keyof Data
        selected: string[]
        data: any
    };

    constructor(props:any) {
        super(props);
        this.state = {
            order: "asc",
            orderBy: "IP",
            selected: [],
            data: []
        };
        // this.data = null;
    }

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

    // UNUSED BUT NOT DELETE
    // handleMultiClick = (event: React.MouseEvent<unknown>, IP: string) => {
    //     const selectedIndex = this.state.selected.indexOf(IP);
    //     let newSelected: string[] = [];
    //
    //      if (selectedIndex === -1) { // если его нет в списке
    //         newSelected = newSelected.concat(this.state.selected, IP); // то добавляем в список
    //     } else if (selectedIndex === 0) { // если уже есть в списке и он на первом месте, то выкидываем его из списка
    //         newSelected = newSelected.concat(this.state.selected.slice(1)); // путем переприсваивания массива с 1 элемента
    //     } else if (selectedIndex === this.state.selected.length - 1) { // если он уже есть в списке и на последнем месте, то выкидываем его из списка
    //         newSelected = newSelected.concat(this.state.selected.slice(0, -1)); // путем переприсваивания массива до length-1 элемента
    //     } else if (selectedIndex > 0) { // в противном случае, если он находится в середине, то выкидываем его из списка
    //         newSelected = newSelected.concat( // путем конкатенации массивов
    //             this.state.selected.slice(0, selectedIndex),
    //             this.state.selected.slice(selectedIndex + 1),
    //         );
    //     }
    //
    //     this.setState({
    //         selected: newSelected
    //     });
    // };

    handleClick = (event: React.MouseEvent<unknown>, IP: string) => {
        const selectedIndex = this.state.selected.indexOf(IP);
        let newSelected: string[] = [];
        if (selectedIndex === -1) { // если его нет в списке
            newSelected = newSelected.concat(newSelected, IP); // то добавляем в список
            this.props.setId(IP);
        } else {
            this.props.setId(undefined);
        }

        this.setState({
            selected: newSelected
        });
    };

    isSelected = (IP: string) => this.state.selected.indexOf(IP) !== -1;

    componentDidMount(): void {
        console.log("table mounted");
        this.client = new WebSocket3('ws://localhost/update/hosts');
        this.client.onmessage = (message) => {
            let data = JSON.parse(message.data);

            this.state.data = [];

            data.forEach((elem, idx, array) => {
                this.state.data.push(createData(elem.Id, elem.Name, elem.IP, elem.Status));
            });

            this.forceUpdate();
            this.client.send("ok");
            // setTimeout(() => {this.client.send("ok")}, 1000);
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
    }

    componentWillUnmount(): void {
        console.log("table unmounted");
        this.client.close(1000, "finish update")
    }

    humanReadStatus(variable:any): any {
        switch (variable) {
            case 0:
                return "Нет соединения";
            case 1:
                return "Соединение установлено";
            case 2:
                return "Соединение потеряно";
            default:
                return "Неопределено"
        }
    }

    render(){
    return (
        <Box sx={{ width: '100%', height: '100%'}}>
                <TableContainer component={Paper} sx={{
                    width: '100%',
                    minWidth:'100%',
                    bgcolor: 'background.paper',
                    position: 'relative',
                    overflow: 'auto',
                    height:'100%',
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
                        />
                        <TableBody>
                            {this.state.data.sort(getComparator(this.state.order, this.state.orderBy))
                                .map((row, index) => {
                                    const isItemSelected = this.isSelected(row.IP);
                                    const labelId = `enhanced-table-checkbox-${index}`;

                                    return (
                                        <TableRow
                                            hover
                                            onClick={(event) => this.handleClick(event, row.IP)}
                                            role="checkbox"
                                            aria-checked={isItemSelected}
                                            tabIndex={-1}
                                            key={row.name}
                                            selected={isItemSelected}
                                        >
                                            <TableCell padding="checkbox">
                                                <Checkbox
                                                    color="primary"
                                                    checked={isItemSelected}
                                                    inputProps={{
                                                        'aria-labelledby': labelId,
                                                    }}
                                                />
                                            </TableCell>
                                            <TableCell
                                                id={labelId}
                                                scope="row"
                                                padding="none"
                                            >
                                                {row.IP}
                                            </TableCell>
                                            <TableCell align="right">{row.name}</TableCell>
                                            <TableCell align="right">{this.humanReadStatus(row.status)}</TableCell>
                                        </TableRow>
                                    );
                                })}

                        </TableBody>
                    </Table>
                </TableContainer>
        </Box>
    )};
}

export default connect(null, mapDispatchToProps)(TableConnections);