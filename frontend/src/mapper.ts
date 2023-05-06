import { RuleStatus } from "./types";

const statusToColor: { [status in RuleStatus]: string } = {
    skip: '#999999',
    pass: '#43A047',
    warn: '#FB8C00',
    fail: '#EF5350',
    error: '#F44336'
}

export const mapStatus = (status: RuleStatus): string => statusToColor[status] || statusToColor['skip']