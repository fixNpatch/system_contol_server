import * as React from 'react';
import Plot from 'react-plotly.js';
import {Box} from "@mui/material";
import {width} from "@mui/system";


export default class PlotComponent extends React.Component {

    getData () {
        return Math.random()
    }

    render() {
        return (
            <Plot style={{height:"100%", width:"90%", margin:"0 auto"}}
                data={[
                    {
                        // x: [1, 2, 3],
                        // y: [this.getData()],

                        x: [1, 2, 3],
                        y: [1, 2, 3],
                        type: 'scatter',
                        fill: 'tozeroy',
                        marker: {color: 'red'},
                    },

                ]}
                layout={{title:"Двойное нажатие возвращает к изначальному состоянию", autosize: true}}
                config={
                    {
                        displayModeBar: false, // this is the line that hides the bar.
                        displaylogo: false, // hide logo
                        // responsive: true
                        // autosizable: false
                    }
                }
            />
        );
    }
}