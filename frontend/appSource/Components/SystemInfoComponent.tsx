import * as React from 'react';
import {List, ListItem, ListItemText, ListSubheader} from "@mui/material";
import { w3cwebsocket as WebSocket3 } from "websocket";
import {connect} from "react-redux";
import store from "../redux/store";

const mapStateToProps = (state:any) => {
    return {
        observedId: state
    };
};

interface Data {
    CpuCores: number;
    CpuModel: string;
    CpuUsage: number[];

    DiskTotal: number;
    DiskFree: number;
    DiskUsed: number;
    DiskPercent: number;

    HostOs: string;
    HostPlatform: string;
    HostPlatformVersion: string;
    HostProcs: number;

}

function createData(
    CpuCores: number,
    CpuModel: string,
    CpuUsage: number[],
    DiskTotal: number,
    DiskFree: number,
    DiskUsed: number,
    DiskPercent: number,
    HostOs: string,
    HostPlatform: string,
    HostPlatformVersion: string,
    HostProcs: number,
): Data {
    return {
        CpuCores,
        CpuModel,
        CpuUsage,
        DiskTotal,
        DiskFree,
        DiskUsed,
        DiskPercent,
        HostOs,
        HostPlatform,
        HostPlatformVersion,
        HostProcs,
    };
}


class SystemInfoComponent extends React.Component<any, any> {
    private client: any;
    state: {
        sav: any
        data: Data
    };

    constructor(props:any) {
        super(props);
        this.state = {
            sav: null,
            data: new class implements Data {
                DiskFree: number;
                DiskPercent: number;
                DiskUsed: number;
                HostOs: string;
                HostPlatform: string;
                HostPlatformVersion: string;
                HostProcs: number;
                CpuCores: number;
                CpuModel: string;
                CpuUsage: number[];
                DiskTotal: number;
            },
        };
    }

    componentWillUnmount = ():void => {
        console.log("unmount");
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

    createConnection = (id:any):void => {
        this.setState({
            sav: this.props.observedId
        });
        this.client = new WebSocket3('ws://localhost/update/info?ip=' + id);
        this.client.onmessage = (message) => {
            let data = JSON.parse(message.data);
            console.log(data);

            this.state.data = createData(
                data.Cpu.Cores,
                data.Cpu.Model,
                data.Cpu.Percentage[0],
                data.Disk.Total,
                data.Disk.Free,
                data.Disk.Used,
                data.Disk.UsedPercent,
                data.Host.OS,
                data.Host.Platform,
                data.Host.HostPlatformVersion,
                data.Host.Procs);

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

    closeConnection = ():void => {
        if (this.client != null) {
            this.client.close(1000, "finish update host info");
        }
        this.client = null;
    };

    render() {
        return (

            <List
                sx={{
                    width: '100%',
                    bgcolor: 'background.paper',
                    position: 'relative',
                    overflow: 'auto',
                    height: "100%",
                    p: 0, // padding
                    m: 0, // margin
                    '& ul': { padding: 0 },
                }}
                subheader={<li/>}
            >

                <li key={`host`}>
                    <ul>
                        <ListSubheader sx={{background:"#565f6e", color:"white"}}>{`Информация о системе`}</ListSubheader>
                        <ListItem key={`host-os`}>
                            <ListItemText primary={`ОС`}/><div>{this.state.data.HostOs}</div>
                        </ListItem>
                        <ListItem key={`host-platform`}>
                            <ListItemText primary={`Платформа`} /><div>{this.state.data.HostPlatform}</div>
                            {/* Example */}
                            {/*<ListItemText sx={{fontSize:15}} primary={`Platform`} disableTypography /> */}
                        </ListItem>
                        <ListItem key={`host-platformVersion`}>
                            <ListItemText primary={`Версия`} /><div>{this.state.data.HostPlatformVersion}</div>
                        </ListItem>
                    </ul>
                </li>

                <li key={`cpu`}>
                    <ul>
                        <ListSubheader sx={{background:"#565f6e", color:"white"}}>{`Информация о процессоре`}</ListSubheader>
                        <ListItem key={`cpu-model`}>
                            <ListItemText primary={`Модель`} /><div>{this.state.data.CpuModel}</div>
                        </ListItem>
                        <ListItem key={`cpu-cores`}>
                            <ListItemText primary={`Количество ядер`} /><div>{this.state.data.CpuCores}</div>
                        </ListItem>
                        <ListItem key={`cpu-percentage`}>
                            <ListItemText primary={`Нагрузка %`} /><div>{this.state.data.CpuUsage}</div>
                        </ListItem>
                    </ul>
                </li>

                <li key={`disk`}>
                    <ul>
                        <ListSubheader sx={{background:"#565f6e", color:"white"}}>{`Дисковое пространство`}</ListSubheader>
                        <ListItem key={`disk-total`}>
                            <ListItemText primary={`Всего`} /><div>{this.state.data.DiskTotal}</div>
                        </ListItem>
                        <ListItem key={`disk-free`}>
                            <ListItemText primary={`Свободно`} /><div>{this.state.data.DiskFree}</div>
                        </ListItem>
                        <ListItem key={`disk-used`}>
                            <ListItemText primary={`Занято`} /><div>{this.state.data.DiskUsed}</div>
                        </ListItem>
                        <ListItem key={`disk-percentage`}>
                            <ListItemText primary={`Использовано %`} /><div>{this.state.data.DiskPercent}</div>
                        </ListItem>
                    </ul>
                </li>
            </List>

        );
    }
}

export default connect(mapStateToProps)(SystemInfoComponent);