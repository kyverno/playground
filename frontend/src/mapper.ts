import { layoutTheme } from './config'
import type { RuleStatus } from './types'

const statusToColor: { [status in RuleStatus]: string } = {
  skip: '#616161',
  pass: '#43A047',
  warn: '#FB8C00',
  fail: '#EF5350',
  error: '#F44336',
  'no match': '#000000'
}

const statusToDarkColor: { [status in RuleStatus]: string } = {
  skip: '#636363',
  pass: '#1B5E20',
  warn: '#FF6F00',
  fail: '#D32F2F',
  error: '#B71C1C',
  'no match': '#000000'
}

export const mapStatus = (status: RuleStatus): string => {
  if (layoutTheme.value === 'dark') {
    return statusToDarkColor[status] || statusToDarkColor['skip']
  }

  return statusToColor[status] || statusToColor['skip']
}
