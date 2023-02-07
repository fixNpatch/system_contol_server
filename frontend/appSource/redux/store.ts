import { createStore } from "redux";
import rootReducer from "./reducers";
import connections from "./reducers/connections";

// export default createStore(rootReducer);

export default createStore(connections);