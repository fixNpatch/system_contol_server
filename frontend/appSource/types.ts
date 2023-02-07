interface HeadCell {
    disablePadding: boolean;
    id: keyof Data;
    label: string;
    numeric: boolean;
}

interface Data {
    id: number;
    name: string;
    IP: string;
    status: number;
}

type Order = 'asc' | 'desc';