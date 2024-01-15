import isObject from 'lodash/isObject';

interface AnyObject {
    [x: string]: any;
}
type SetValue = string | AnyObject;
type GetResponse = string | AnyObject | null;

const getParser = (value: string | null): GetResponse => {
    if (value === null) {
        return null;
    }
    try {
        return JSON.parse(value);
    } catch (error) {
        return value;
    }
};

const setParser = (value: SetValue): string => {
    let val;
    if (isObject(value)) {
        val = JSON.stringify(value);
    } else {
        val = value;
    }
    return val;
};

export const lsGetItem = (key: string): GetResponse => {
    const val = localStorage.getItem(key);
    return getParser(val);
};

export const lsSetItem = (key: string, value: SetValue): void => {
    const val = setParser(value);
    localStorage.setItem(key, val);
};

export const getLocalStorage = {
    getItem: lsGetItem,
    setItem: lsSetItem,
};
