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
import {visuallyHidden} from '@mui/utils';
import {connect} from "react-redux";
import * as $ from 'jquery';
import {setObservedConnection} from "../redux/actions";
import {getComparator} from "../comparators";
import BuildIcon from '@mui/icons-material/Build';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';


function createData(
    id: number,
    name: string,
    IP: string,
): any {
    return {
        id,
        name,
        IP,
    };
}


const headCells: readonly HeadCell[] = [
    {
        id: 'id',
        numeric: false,
        disablePadding: false,
        label: 'ID',
    },
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
    }

];

interface EnhancedTableProps {
    onRequestSort: (event: React.MouseEvent<unknown>, property: keyof Data) => void;
    order: Order;
    orderBy: string;
}


function EnhancedTableHead(props: EnhancedTableProps) {
    const {order, orderBy, onRequestSort} =
        props;
    const createSortHandler =
        (property: keyof Data) => (event: React.MouseEvent<unknown>) => {
            onRequestSort(event, property);
        };

    return (
        <TableHead>
            <TableRow>
                <TableCell/>
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


class ManagementTableConnections extends React.Component<any, any> {
    // private data:any;
    state: {
        order: Order
        orderBy: keyof Data
        data: any
    };
    private client: any;

    constructor(props: any) {
        super(props);
        this.state = {
            order: "asc",
            orderBy: "IP",
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


    componentDidMount = () => {

        this.setState({
            data:[]
        });

        console.log("table mounted");
        $.get("/management/list").done((data:any) => {
            console.log("success");
            console.log(data);

            this.state.data = [];

            data.forEach((elem, idx, array) => {
                this.state.data.push(createData(elem.Id, elem.Name, elem.IP));
            });

            this.forceUpdate();
        }).fail((): any => {
            console.log("cannot get management/list");
        });
    };

    componentWillUnmount(): void {
        console.log("table unmounted");
    }


    render() {
        return (
            <Box sx={{width: '100%', height: '100%'}}>
                <TableContainer component={Paper} sx={{
                    width: '100%',
                    minWidth: '100%',
                    bgcolor: 'background.paper',
                    position: 'relative',
                    overflow: 'auto',
                    height: '100%',
                    fontSize: 10,
                    '& ul': {padding: 0},
                }}>
                    <Table
                        sx={{
                            minWidth: 650, "& th": {
                                // color: "rgba(96, 96, 96)",
                                backgroundColor: '#565f6e',
                                color: "white"
                            }
                        }} size="small" aria-label="List of followed hosts" stickyHeader
                    >
                        <EnhancedTableHead
                            order={this.state.order}
                            orderBy={this.state.orderBy}
                            onRequestSort={this.handleRequestSort}
                        />
                        <TableBody>
                            {this.state.data.sort(getComparator(this.state.order, this.state.orderBy))
                                .map((row, index) => {
                                    return (
                                        <TableRow
                                            hover
                                            role="checkbox"
                                            tabIndex={-1}
                                            key={row.id}
                                        >
                                            <TableCell>
                                                <BuildIcon/>
                                            </TableCell>
                                            <TableCell>
                                                <DeleteForeverIcon/>
                                            </TableCell>
                                            <TableCell>
                                                {row.id}
                                            </TableCell>
                                            <TableCell
                                                scope="row"
                                                padding="none"
                                            >
                                                {row.IP}
                                            </TableCell>
                                            <TableCell align="right">{row.name}</TableCell>
                                        </TableRow>
                                    );
                                })}

                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        )
    };
}

export default ManagementTableConnections;