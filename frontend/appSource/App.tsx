import * as React from 'react';
import {ThemeProvider} from "@mui/material/styles";
import theme from "./Components/theme";
import MiniDrawer from "./Components/MiniDrawer";
import {CssBaseline} from "@mui/material";
import {Provider} from "react-redux";
import store from "./redux/store";


export default class App extends React.Component{
    constructor(props:any) {
        super(props);
        this.state = {controlPanelOpen: false};
    }
    render() {
        return (
            <ThemeProvider theme={theme}>
                <Provider store={store}>
                    <CssBaseline/>
                    <MiniDrawer/>
                </Provider>
            </ThemeProvider>
        );
    }
}