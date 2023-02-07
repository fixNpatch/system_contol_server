import {SET_OBSERVED_CONNECTION} from "../actionTypes";

const initialState = {
    observedId: undefined
};

export default function(state = initialState.observedId, action) {
    switch (action.type) {
        case SET_OBSERVED_CONNECTION: {
            return action.payload;
        }
        default:
            return state;
    }
}