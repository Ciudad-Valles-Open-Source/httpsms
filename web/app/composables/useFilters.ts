import {
  formatPhoneNumber,
  phoneCountry,
  formatTimestamp,
  formatMoney,
  formatDecimal,
  formatBillingPeriod,
  formatBillingPeriodDateOrdinal,
  humanizeTime,
  humanizeTimeShort,
  startsWithLetter,
} from '../utils/filters'
import { capitalize } from '../utils/capitalize'

// Dates, durations and numbers follow the app language, not the browser's.
export function useFilters() {
  const locale = () => useNuxtApp().$i18n.locale.value
  return {
    formatPhoneNumber,
    phoneCountry,
    formatTimestamp: (value: string) => formatTimestamp(value, locale()),
    formatMoney,
    formatDecimal: (value: string | number) => formatDecimal(value, locale()),
    formatBillingPeriod,
    humanizeTimeShort: (date: string) => humanizeTimeShort(date, locale()),
    formatBillingPeriodDateOrdinal,
    humanizeTime: (value: string) => humanizeTime(value, locale()),
    startsWithLetter,
    capitalize,
  }
}
