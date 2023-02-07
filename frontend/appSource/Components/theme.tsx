import {createTheme} from "@mui/material/styles";

// declare module '@mui/material/styles' {
//     interface Theme {
//         status: {
//             danger: string;
//         };
//     }
//     // allow configuration using `createTheme`
//     interface ThemeOptions {
//         status?: {
//             danger?: string;
//         };
//     }
// }

// A custom theme for this app
const theme = createTheme({
    typography: {
        fontSize: 18,
    },
    palette: {
        primary: {
            main: '#565f6e',
        },
        secondary: {
            main: '#19857b',
        },
        error: {
            main: "#6e1214",
        },
        background: {
            default: '#dbe1f1',
        },
    },
});

export default theme;