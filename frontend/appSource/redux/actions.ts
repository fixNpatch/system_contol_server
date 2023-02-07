import {SET_OBSERVED_CONNECTION} from "./actionTypes";

export const setObservedConnection = (id:any) => ({
    type: SET_OBSERVED_CONNECTION,
    payload: id
});