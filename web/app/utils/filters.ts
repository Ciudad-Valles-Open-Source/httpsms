import {
  intervalToDuration,
  formatDuration,
  differenceInSeconds,
  differenceInMinutes,
  differenceInDays,
  differenceInHours,
  type Locale,
} from 'date-fns'
import { enUS, es } from 'date-fns/locale'
import { parsePhoneNumber, isValidPhoneNumber } from 'libphonenumber-js'

export function formatPhoneNumber(value: string): string {
  if (!value || typeof value !== 'string') {
    return value ?? ''
  }
  if (!isValidPhoneNumber(value)) {
    return value
  }
  const phoneNumber = parsePhoneNumber(value)
  if (phoneNumber) {
    return phoneNumber.formatInternational()
  }
  return value
}

export function phoneCountry(value: string): string {
  const phoneNumber = parsePhoneNumber(value)
  if (phoneNumber && phoneNumber.country) {
    const regionNames = new Intl.DisplayNames(undefined, { type: 'region' })
    return regionNames.of(phoneNumber.country) ?? 'Earth'
  }
  return 'Earth'
}

/** date-fns locale for an app locale code (`es`, `en`). */
export function dateFnsLocale(locale = 'en'): Locale {
  return locale.startsWith('es') ? es : enUS
}

export function formatTimestamp(value: string, locale?: string): string {
  return new Date(value).toLocaleString(locale)
}

export function formatMoney(value: string | number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(typeof value === 'string' ? parseInt(value) : value)
}

export function formatDecimal(
  value: string | number,
  locale = 'en-US',
): string {
  return new Intl.NumberFormat(locale, {
    style: 'decimal',
  }).format(typeof value === 'string' ? parseInt(value) : value)
}

export function formatBillingPeriod(value: string): string {
  return new Date(value).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
  })
}

export function formatBillingPeriodDateOrdinal(value: string): string {
  const date = new Date(value)
  const day = date.getDate()
  const month = date.toLocaleDateString('en-US', { month: 'long' })
  const year = date.getFullYear()

  const suffix =
    day % 10 === 1 && day !== 11
      ? 'st'
      : day % 10 === 2 && day !== 12
        ? 'nd'
        : day % 10 === 3 && day !== 13
          ? 'rd'
          : 'th'

  return `${month} ${day}<sup>${suffix}</sup> ${year}`
}

export interface BillingPeriodDateOrdinalParts {
  leading: string
  suffix: string
  trailing: string
}

export function formatBillingPeriodDateOrdinalParts(
  value: string,
): BillingPeriodDateOrdinalParts {
  const date = new Date(value)
  const day = date.getDate()
  const month = date.toLocaleDateString('en-US', { month: 'long' })
  const year = date.getFullYear()

  const suffix =
    day % 10 === 1 && day !== 11
      ? 'st'
      : day % 10 === 2 && day !== 12
        ? 'nd'
        : day % 10 === 3 && day !== 13
          ? 'rd'
          : 'th'

  return { leading: `${month} ${day}`, suffix, trailing: ` ${year}` }
}

export function humanizeTime(value: string, locale = 'en'): string {
  const durations = intervalToDuration({
    start: new Date(value),
    end: new Date(),
  })
  return formatDuration(durations, { locale: dateFnsLocale(locale) })
}

export function startsWithLetter(value: string): boolean {
  return /^[a-zA-Z]/.test(value)
}

/** "2 minutes ago" / "hace 2 minutos", in the given locale. */
export function humanizeTimeShort(date: string, locale = 'en'): string {
  const now = new Date()
  const then = new Date(date)
  const relative = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
  const seconds = differenceInSeconds(now, then)
  const minutes = differenceInMinutes(now, then)
  const hours = differenceInHours(now, then)
  const days = differenceInDays(now, then)

  if (seconds < 60) return relative.format(0, 'second')
  if (minutes < 60) return relative.format(-minutes, 'minute')
  if (hours < 24) return relative.format(-hours, 'hour')
  return relative.format(-days, 'day')
}
